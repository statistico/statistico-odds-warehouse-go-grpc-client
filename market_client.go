package statisticooddswarehouse

import (
	"context"
	"github.com/statistico/statistico-proto/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

type MarketClient interface {
	MarketRunnerSearch(ctx context.Context, r *statisticoproto.MarketRunnerRequest) (<-chan *statisticoproto.MarketRunner, error)
}

type marketClient struct {
	client statisticoproto.MarketServiceClient
}

func (m *marketClient) MarketRunnerSearch(ctx context.Context, r *statisticoproto.MarketRunnerRequest) (<-chan *statisticoproto.MarketRunner, error) {
	runners := make(chan *statisticoproto.MarketRunner, 100)
	defer close(runners)

	stream, err := m.client.MarketRunnerSearch(ctx, r)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				return runners, ErrorInvalidArgument{err}
			case codes.Internal:
				return runners, ErrorInternalServerError{err}
			default:
				return runners, ErrorBadGateway{err}
			}
		}

		return runners, err
	}

	for {
		st, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			return runners, ErrorInternalServerError{err}
		}

		runners <- st
	}

	return runners, nil
}

func NewMarketClient(p statisticoproto.MarketServiceClient) MarketClient {
	return &marketClient{client: p}
}
