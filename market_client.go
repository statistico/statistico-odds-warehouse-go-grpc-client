package statisticooddswarehouse

import (
	"context"
	"github.com/statistico/statistico-proto/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

type MarketClient interface {
	MarketRunnerSearch(ctx context.Context, r *statisticoproto.MarketRunnerRequest, errCh chan error) <-chan *statisticoproto.MarketRunner
}

type marketClient struct {
	client statisticoproto.MarketServiceClient
}

func (m *marketClient) MarketRunnerSearch(ctx context.Context, r *statisticoproto.MarketRunnerRequest, errCh chan error) <-chan *statisticoproto.MarketRunner {
	runners := make(chan *statisticoproto.MarketRunner, 100)

	stream, err := m.client.MarketRunnerSearch(ctx, r)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				return sendError(ErrorInvalidArgument{err}, errCh, runners)
			case codes.Internal:
				return sendError(ErrorInternalServerError{err}, errCh, runners)
			default:
				return sendError(ErrorBadGateway{err}, errCh, runners)
			}
		}

		return sendError(err, errCh, runners)
	}

	go streamMarketRunners(stream, runners, errCh)

	return runners
}

func streamMarketRunners(stream statisticoproto.MarketService_MarketRunnerSearchClient, ch chan *statisticoproto.MarketRunner, errCh chan error) {
	for {
		st, err := stream.Recv()

		if err == io.EOF {
			close(ch)
			break
		}

		if err != nil {
			errCh <- ErrorInternalServerError{err}
			close(ch)
			break
		}

		ch <- st
	}
}

func sendError(err error, errCh chan error, ch chan *statisticoproto.MarketRunner) <-chan *statisticoproto.MarketRunner {
	errCh <- err
	close(ch)
	return ch
}

func NewMarketClient(p statisticoproto.MarketServiceClient) MarketClient {
	return &marketClient{client: p}
}
