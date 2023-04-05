package statisticooddswarehouse

import (
	"context"
	"github.com/statistico/statistico-proto/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

type MarketClient interface {
	GetExchangeOdds(ctx context.Context, r *statistico.ExchangeOddsRequest) ([]*statistico.ExchangeOdds, error)
}

type marketClient struct {
	client statistico.OddsWarehouseServiceClient
}

func (m *marketClient) GetExchangeOdds(ctx context.Context, r *statistico.ExchangeOddsRequest) ([]*statistico.ExchangeOdds, error) {
	odds := []*statistico.ExchangeOdds{}

	stream, err := m.client.GetExchangeOdds(ctx, r)

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.Internal:
				return odds, ErrorInternalServerError{err}
			default:
				return odds, ErrorBadGateway{err}
			}
		}

		return odds, err
	}

	for {
		o, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			return odds, ErrorInternalServerError{err}
		}

		odds = append(odds, o)
	}

	return odds, nil
}

func NewMarketClient(p statistico.OddsWarehouseServiceClient) MarketClient {
	return &marketClient{client: p}
}
