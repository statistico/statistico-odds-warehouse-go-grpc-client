package statisticooddswarehouse

import (
	"context"
	"github.com/statistico/statistico-proto/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

const DefaultChannelSize = 5000

type MarketClient interface {
	MarketRunnerSearch(ctx context.Context, r *statistico.MarketRunnerRequest, chSize int) (<-chan *statistico.MarketRunner, <-chan error)
}

type marketClient struct {
	client statistico.OddsWarehouseServiceClient
}

func (m *marketClient) MarketRunnerSearch(ctx context.Context, r *statistico.MarketRunnerRequest, chSize int) (<-chan *statistico.MarketRunner, <-chan error) {
	if chSize == 0 {
		chSize = DefaultChannelSize
	}

	runners := make(chan *statistico.MarketRunner, chSize)
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

func streamMarketRunners(stream statistico.OddsWarehouseService_MarketRunnerSearchClient, ch chan *statistico.MarketRunner, errCh chan error) {
	for {
		st, err := stream.Recv()

		if err == io.EOF {
			close(ch)
			close(errCh)
			break
		}

		if err != nil {
			errCh <- ErrorStreamFailure{item: st, err: err}
			break
		}

		ch <- st
	}
}

func sendError(err error, ch chan *statistico.MarketRunner, errCh chan error) (<-chan *statistico.MarketRunner, <-chan error) {
	errCh <- err
	close(ch)
	close(errCh)
	return ch, errCh
}

func NewMarketClient(p statistico.OddsWarehouseServiceClient) MarketClient {
	return &marketClient{client: p}
}
