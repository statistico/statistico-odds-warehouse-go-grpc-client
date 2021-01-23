# Statistico Odds Warehouse Go gRPC Client

[![CircleCI](https://circleci.com/gh/statistico/statistico-odds-warehouse-go-grpc-client/tree/main.svg?style=shield)](https://circleci.com/gh/statistico/statistico-betfair-go-client/tree/master)

This library is a Go client for the Statistico Odds Warehouse service. API reference can be found here:

[Statistico Odds Warehouse Proto](https://github.com/statistico/statistico-proto/blob/main/market.proto)

## Installation
```.env
$ go get -u github.com/statistico/statistico-odds-warehouse-go-grpc-client
```
## Usage
To instantiate the required client struct and retrieve and search for marker runner resources:
```go
package main

import (
    "context"
    "github.com/golang/protobuf/ptypes/timestamp"
    "github.com/statistico/statistico-odds-warehouse-go-grpc-client"
    "github.com/statistico/statistico-proto/go"
    "google.golang.org/grpc"
)

func main() {
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    
    if err != nil {
        // Handle error
    }

    c := statistico.NewMarketServiceClient(conn)

    client := statisticooddswarehouse.NewMarketClient(c)

    request := statistico.MarketRunnerRequest{
        Name:                 "MATCH_ODDS",
        RunnerFilter:         &statistico.RunnerFilter{
            Name:                 "Home",
            Line:                 statistico.LineEnum_CLOSING,
            Operators:            []*statistico.MetricOperator{
                    {
                        Metric: statistico.MetricEnum_GTE,
                        Value: 2.50,
                    },
                },
        },
        CompetitionIds:       []uint64{1, 2, 3},
        SeasonIds:            []uint64{4, 5, 6},
        DateFrom:             &timestamp.Timestamp{Seconds: 1584014400},
        DateTo:               &timestamp.Timestamp{Seconds: 1584014400},
    }

    ctx := context.Background()
    errCh := make(chan error, 1)

    markets := client.MarketRunnerSearch(ctx, &request, errCh)
    
    select {
    case market := <-markets:
        // Handle with market
    case err := <-errCh:
        // Handle err
    }
}
```