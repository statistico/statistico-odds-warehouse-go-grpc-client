package statisticooddswarehouse

import (
	"context"
	"github.com/statistico/statistico-proto/statistico-odds-warehouse/go"
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
	ch := make(chan *statisticoproto.MarketRunner, 100)
	errCh := make(chan error)

	go m.streamMarketRunners(ctx, r, ch, errCh)

	return ch, errCh
}

func (m *marketClient) streamMarketRunners(ctx context.Context, r *statisticoproto.MarketRunnerRequest, ch chan<- *statisticoproto.MarketRunner, errCh chan<- error) {
	stream, err := m.client.MarketRunnerSearch(ctx, r)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				errCh <- ErrorInvalidArgument{err}
				break
			case codes.Internal:
				errCh <- ErrorInternalServerError{err}
				break
			default:
				errCh <- ErrorBadGateway{err}
				break
			}
		}

		closeChannels(ch, errCh)
		return
	}

	for {
		mr, err := stream.Recv()

		if err == io.EOF {
			closeChannels(ch, errCh)
			return
		}

		if err != nil {
			errCh <- ErrorInternalServerError{err}
			closeChannels(ch, errCh)
			return
		}

		ch <- mr
	}
}

func closeChannels(ch chan<- *statisticoproto.MarketRunner, errCh chan<- error) {
	close(ch)
	close(errCh)
}

func NewMarketClient(p statisticoproto.MarketServiceClient) MarketClient {
	return &marketClient{client: p}
}
