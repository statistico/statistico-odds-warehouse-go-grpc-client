package statisticooddswarehouse_test

import (
	"context"
	"errors"
	"github.com/statistico/statistico-odds-warehouse-go-grpc-client"
	"github.com/statistico/statistico-proto/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"testing"
)

func TestMarketClient_MarketRunnerSearch(t *testing.T) {
	t.Run("calls market client and returns a slice of ExchangeOdds struct", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoMarketClient)
		client := statisticooddswarehouse.NewMarketClient(m)

		stream := new(MockMarketStream)

		request := statistico.ExchangeOddsRequest{
			EventId:  1234,
			Market:   "OVER_UNDER_25",
			Exchange: "BETFAIR",
			Runner:   "OVER",
			Limit:    3,
		}

		o1 := statistico.ExchangeOdds{Price: 1.95, Timestamp: 1680706966}
		o2 := statistico.ExchangeOdds{Price: 1.97, Timestamp: 1680706956}
		o3 := statistico.ExchangeOdds{Price: 1.98, Timestamp: 1680706936}

		ctx := context.Background()

		m.On("GetExchangeOdds", ctx, &request).Return(stream, nil)
		stream.On("Recv").Once().Return(&o1, nil)
		stream.On("Recv").Once().Return(&o2, nil)
		stream.On("Recv").Once().Return(&o3, nil)
		stream.On("Recv").Once().Return(&statistico.ExchangeOdds{}, io.EOF)

		odds, err := client.GetExchangeOdds(ctx, &request)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		a := assert.New(t)
		a.Equal(float32(1.95), odds[0].Price)
		a.Equal(uint64(1680706966), odds[0].Timestamp)
		a.Equal(float32(1.97), odds[1].Price)
		a.Equal(uint64(1680706956), odds[1].Timestamp)
		a.Equal(float32(1.98), odds[2].Price)
		a.Equal(uint64(1680706936), odds[2].Timestamp)
		m.AssertExpectations(t)
		stream.AssertExpectations(t)
	})

	t.Run("returns internal server error", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoMarketClient)
		client := statisticooddswarehouse.NewMarketClient(m)

		stream := new(MockMarketStream)

		request := statistico.ExchangeOddsRequest{
			EventId:  1234,
			Market:   "OVER_UNDER_25",
			Exchange: "BETFAIR",
			Runner:   "OVER",
			Limit:    3,
		}

		ctx := context.Background()

		e := status.Error(codes.Internal, "internal error")

		m.On("GetExchangeOdds", ctx, &request).Return(stream, e)

		_, err := client.GetExchangeOdds(ctx, &request)

		if err == nil {
			t.Fatal("Expected errors, got nil")
		}

		assert.Equal(t, "internal server error returned from external service: rpc error: code = Internal desc = internal error", err.Error())
		m.AssertExpectations(t)
		stream.AssertExpectations(t)
	})

	t.Run("returns internal server error if error parsing stream", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoMarketClient)
		client := statisticooddswarehouse.NewMarketClient(m)

		stream := new(MockMarketStream)

		request := statistico.ExchangeOddsRequest{
			EventId:  1234,
			Market:   "OVER_UNDER_25",
			Exchange: "BETFAIR",
			Runner:   "OVER",
			Limit:    3,
		}

		ctx := context.Background()

		e := errors.New("oh damn")

		o1 := statistico.ExchangeOdds{Price: 1.95, Timestamp: 1680706966}
		o2 := statistico.ExchangeOdds{Price: 1.97, Timestamp: 1680706956}
		o3 := statistico.ExchangeOdds{Price: 1.98, Timestamp: 1680706936}

		m.On("GetExchangeOdds", ctx, &request).Return(stream, nil)
		stream.On("Recv").Once().Return(&o1, nil)
		stream.On("Recv").Once().Return(&o2, nil)
		stream.On("Recv").Once().Return(&o3, nil)
		stream.On("Recv").Once().Return(&statistico.ExchangeOdds{}, e)

		_, err := client.GetExchangeOdds(ctx, &request)

		if err == nil {
			t.Fatal("Expected errors, got nil")
		}

		assert.Equal(t, "internal server error returned from external service: oh damn", err.Error())
		m.AssertExpectations(t)
		stream.AssertExpectations(t)
	})
}

type MockProtoMarketClient struct {
	mock.Mock
}

func (m *MockProtoMarketClient) GetExchangeOdds(ctx context.Context, r *statistico.ExchangeOddsRequest, opts ...grpc.CallOption) (statistico.OddsWarehouseService_GetExchangeOddsClient, error) {
	args := m.Called(ctx, r)
	return args.Get(0).(statistico.OddsWarehouseService_GetExchangeOddsClient), args.Error(1)
}

type MockMarketStream struct {
	mock.Mock
	grpc.ClientStream
}

func (m *MockMarketStream) Recv() (*statistico.ExchangeOdds, error) {
	args := m.Called()
	return args.Get(0).(*statistico.ExchangeOdds), args.Error(1)
}
