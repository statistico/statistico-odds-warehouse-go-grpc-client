package statisticooddswarehouse

import (
	"context"
	"github.com/statistico/statistico-proto/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

type MarketClient interface {
	MarketRunnerSearch(ctx context.Context, r *statisticoproto.MarketRunnerRequest) ([]*statisticoproto.MarketRunner, error)
}

type marketClient struct {
	client statisticoproto.MarketServiceClient
}

func (m *marketClient) MarketRunnerSearch(ctx context.Context, r *statisticoproto.MarketRunnerRequest) ([]*statisticoproto.MarketRunner, error) {
	mr := []*statisticoproto.MarketRunner{}

	stream, err := m.client.MarketRunnerSearch(ctx, r)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.InvalidArgument:
				return mr, ErrorInvalidArgument{err}
			case codes.Internal:
				return mr, ErrorInternalServerError{err}
			default:
				return mr, ErrorBadGateway{err}
			}
		}

		return mr, err
	}

	for {
		st, err := stream.Recv()

		if err == io.EOF {
			return mr, nil
		}

		if err != nil {
			return mr, ErrorInternalServerError{err}
		}

		mr = append(mr, st)
	}
}

func NewMarketClient(p statisticoproto.MarketServiceClient) MarketClient {
	return &marketClient{client: p}
}
