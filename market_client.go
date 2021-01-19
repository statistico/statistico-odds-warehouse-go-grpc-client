package statisticooddswarehouse

import (
	"context"
	"github.com/statistico/statistico-proto/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

type MarketClient interface {
	MarketRunnerSearch(ctx context.Context, r *statisticoproto.MarketRunnerRequest) (<-chan *statisticoproto.MarketRunner, <-chan error)
}

type marketClient struct {
	client statisticoproto.MarketServiceClient
}

func (m *marketClient) MarketRunnerSearch(ctx context.Context, r *statisticoproto.MarketRunnerRequest) (<-chan *statisticoproto.MarketRunner, <-chan error) {
	runners := make(chan *statisticoproto.MarketRunner, 100)
	errCh := make(chan error, 1)

	stream, err := m.client.MarketRunnerSearch(ctx, r)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				return sendError(ErrorInvalidArgument{err}, runners, errCh)
			case codes.Internal:
				return sendError(ErrorInternalServerError{err}, runners, errCh)
			default:
				return sendError(ErrorBadGateway{err}, runners, errCh)
			}
		}

		return sendError(err, runners, errCh)
	}

	go streamMarketRunners(stream, runners, errCh)

	return runners, errCh
}

func streamMarketRunners(stream statisticoproto.MarketService_MarketRunnerSearchClient, ch chan *statisticoproto.MarketRunner, errCh chan error) {
	for {
		st, err := stream.Recv()

		if err == io.EOF {
			close(ch)
			close(errCh)
			break
		}

		if err != nil {
			errCh <- ErrorInternalServerError{err}
			close(ch)
			close(errCh)
			break
		}

		ch <- st
	}
}

func sendError(err error, ch chan *statisticoproto.MarketRunner, errCh chan error) (<-chan *statisticoproto.MarketRunner, <-chan error) {
	errCh <- err
	close(ch)
	close(errCh)
	return ch, errCh
}

func NewMarketClient(p statisticoproto.MarketServiceClient) MarketClient {
	return &marketClient{client: p}
}
