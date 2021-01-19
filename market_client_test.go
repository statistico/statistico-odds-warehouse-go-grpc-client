package statisticooddswarehouse_test

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/timestamp"
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
	t.Run("calls market client and returns a slice of MarketRunner struct", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoMarketClient)
		client := statisticooddswarehouse.NewMarketClient(m)

		stream := new(MockMarketStream)

		request := statisticoproto.MarketRunnerRequest{
			Name: "MATCH_ODDS",
			RunnerFilter: &statisticoproto.RunnerFilter{
				Name: "Home",
				Line: statisticoproto.LineEnum_CLOSING,
				Operators: []*statisticoproto.MetricOperator{
					{
						Metric: statisticoproto.MetricEnum_LTE,
						Value:  2.50,
					},
				},
			},
			CompetitionIds: []uint64{1, 2, 3},
			SeasonIds:      []uint64{4, 5, 6},
			DateFrom:       &timestamp.Timestamp{Seconds: 1584014400},
			DateTo:         &timestamp.Timestamp{Seconds: 1584014400},
		}

		mk1 := newProtoMarketRunner("1.2371761")
		mk2 := newProtoMarketRunner("1.2371762")

		ctx := context.Background()

		m.On("MarketRunnerSearch", ctx, &request, []grpc.CallOption(nil)).Return(stream, nil)
		stream.On("Recv").Once().Return(mk1, nil)
		stream.On("Recv").Once().Return(mk2, nil)
		stream.On("Recv").Once().Return(&statisticoproto.MarketRunner{}, io.EOF)

		ch, errCh := client.MarketRunnerSearch(ctx, &request)

		err := <-errCh

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		a := assert.New(t)
		a.Equal(mk1, <-ch)
		a.Equal(mk2, <-ch)
		m.AssertExpectations(t)
		stream.AssertExpectations(t)
	})

	t.Run("returns error if invalid argument error returned by result client", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoMarketClient)
		client := statisticooddswarehouse.NewMarketClient(m)

		stream := new(MockMarketStream)

		request := statisticoproto.MarketRunnerRequest{
			Name: "MATCH_ODDS",
			RunnerFilter: &statisticoproto.RunnerFilter{
				Name: "Home",
				Line: statisticoproto.LineEnum_CLOSING,
				Operators: []*statisticoproto.MetricOperator{
					{
						Metric: statisticoproto.MetricEnum_LTE,
						Value:  2.50,
					},
				},
			},
			CompetitionIds: []uint64{1, 2, 3},
			SeasonIds:      []uint64{4, 5, 6},
			DateFrom:       &timestamp.Timestamp{Seconds: 1584014400},
			DateTo:         &timestamp.Timestamp{Seconds: 1584014400},
		}

		ctx := context.Background()

		e := status.Error(codes.InvalidArgument, "incorrect format")

		m.On("MarketRunnerSearch", ctx, &request, []grpc.CallOption(nil)).Return(stream, e)

		ch, errCh := client.MarketRunnerSearch(ctx, &request)

		err := <-errCh

		if err == nil {
			t.Fatal("Expected errors, got nil")
		}

		assert.Equal(t, "invalid argument provided: rpc error: code = InvalidArgument desc = incorrect format", err.Error())
		assert.Equal(t, 0, len(ch))
		m.AssertExpectations(t)
		stream.AssertExpectations(t)
	})

	t.Run("returns internal server error", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoMarketClient)
		client := statisticooddswarehouse.NewMarketClient(m)

		stream := new(MockMarketStream)

		request := statisticoproto.MarketRunnerRequest{
			Name: "MATCH_ODDS",
			RunnerFilter: &statisticoproto.RunnerFilter{
				Name: "Home",
				Line: statisticoproto.LineEnum_CLOSING,
				Operators: []*statisticoproto.MetricOperator{
					{
						Metric: statisticoproto.MetricEnum_LTE,
						Value:  2.50,
					},
				},
			},
			CompetitionIds: []uint64{1, 2, 3},
			SeasonIds:      []uint64{4, 5, 6},
			DateFrom:       &timestamp.Timestamp{Seconds: 1584014400},
			DateTo:         &timestamp.Timestamp{Seconds: 1584014400},
		}

		ctx := context.Background()

		e := status.Error(codes.Internal, "internal error")

		m.On("MarketRunnerSearch", ctx, &request, []grpc.CallOption(nil)).Return(stream, e)

		ch, errCh := client.MarketRunnerSearch(ctx, &request)

		err := <-errCh

		if err == nil {
			t.Fatal("Expected errors, got nil")
		}

		assert.Equal(t, "internal server error returned from external service: rpc error: code = Internal desc = internal error", err.Error())
		assert.Equal(t, 0, len(ch))
		m.AssertExpectations(t)
		stream.AssertExpectations(t)
	})

	t.Run("returns internal server error if error parsing stream", func(t *testing.T) {
		t.Helper()

		m := new(MockProtoMarketClient)
		client := statisticooddswarehouse.NewMarketClient(m)

		stream := new(MockMarketStream)

		request := statisticoproto.MarketRunnerRequest{
			Name: "MATCH_ODDS",
			RunnerFilter: &statisticoproto.RunnerFilter{
				Name: "Home",
				Line: statisticoproto.LineEnum_CLOSING,
				Operators: []*statisticoproto.MetricOperator{
					{
						Metric: statisticoproto.MetricEnum_LTE,
						Value:  2.50,
					},
				},
			},
			CompetitionIds: []uint64{1, 2, 3},
			SeasonIds:      []uint64{4, 5, 6},
			DateFrom:       &timestamp.Timestamp{Seconds: 1584014400},
			DateTo:         &timestamp.Timestamp{Seconds: 1584014400},
		}

		ctx := context.Background()

		e := errors.New("oh damn")

		mk1 := newProtoMarketRunner("1.2371761")
		mk2 := newProtoMarketRunner("1.2371762")

		m.On("MarketRunnerSearch", ctx, &request, []grpc.CallOption(nil)).Return(stream, nil)
		stream.On("Recv").Once().Return(mk1, nil)
		stream.On("Recv").Once().Return(mk2, nil)
		stream.On("Recv").Once().Return(&statisticoproto.MarketRunner{}, e)

		ch, errCh := client.MarketRunnerSearch(ctx, &request)

		err := <-errCh

		if err == nil {
			t.Fatal("Expected errors, got nil")
		}

		assert.Equal(t, "internal server error returned from external service: oh damn", err.Error())
		assert.Equal(t, 2, len(ch))
		m.AssertExpectations(t)
		stream.AssertExpectations(t)
	})
}

func newProtoMarketRunner(marketID string) *statisticoproto.MarketRunner {
	return &statisticoproto.MarketRunner{
		MarketId:      marketID,
		MarketName:    "MATCH_ODDS",
		RunnerName:    "Draw",
		EventId:       1982181,
		CompetitionId: 8,
		SeasonId:      17420,
		EventDate:     &timestamp.Timestamp{Seconds: 1584014400},
		Side:          "BACK",
		Exchange:      "betfair",
		Prices:        []*statisticoproto.Price{},
	}
}

type MockProtoMarketClient struct {
	mock.Mock
}

func (m *MockProtoMarketClient) MarketRunnerSearch(ctx context.Context, in *statisticoproto.MarketRunnerRequest, opts ...grpc.CallOption) (statisticoproto.MarketService_MarketRunnerSearchClient, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(statisticoproto.MarketService_MarketRunnerSearchClient), args.Error(1)
}

type MockMarketStream struct {
	mock.Mock
	grpc.ClientStream
}

func (m *MockMarketStream) Recv() (*statisticoproto.MarketRunner, error) {
	args := m.Called()
	return args.Get(0).(*statisticoproto.MarketRunner), args.Error(1)
}
