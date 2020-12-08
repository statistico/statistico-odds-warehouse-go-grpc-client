# Statistico Odds Warehouse Go gRPC Client

[![CircleCI](https://circleci.com/gh/statistico/statistico-odds-warehouse-go-grpc-client/tree/main.svg?style=shield)](https://circleci.com/gh/statistico/statistico-betfair-go-client/tree/master)

This library is a Go client for the Statistico Odds Warehouse service. API reference can be found here:

[Statistico Odds Warehouse Proto](https://github.com/statistico/statistico-proto/tree/main/statistico-odds-warehouse)

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
    "github.com/statistico/statistico-odds-warehouse-go-grpc-client"
    "github.com/statistico/statistico-proto/statistico-odds-warehouse/go"
    "google.golang.org/grpc"
)

func main() {
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    
    if err != nil {
        // Handle error
    }

    c := statisticoproto.NewMarketServiceClient(conn)

    client := statisticooddswarehouse.NewMarketClient(c)

    request := statisticoproto.MarketRunnerRequest{
        Name:                 "MATCH_ODDS",
        RunnerFilter:         &statisticoproto.RunnerFilter{
        Name:                 "Home",
        Line:                 statisticoproto.RunnerFilter_CLOSING,
        Operators:            []*statisticoproto.FilterOperator{
                {
                    Operator: statisticoproto.FilterOperator_LTE,
                    Value: 2.50,
                },
            },
        },
        CompetitionIds:       []uint64{1, 2, 3},
        SeasonIds:            []uint64{4, 5, 6},
        DateFrom:             &wrappers.StringValue{Value: "2020-12-07T12:00:00+00:00"},
        DateTo:               &wrappers.StringValue{Value: "2020-12-07T20:00:00+00:00"},
    }

    ctx := context.Background()

    markets, errCh := client.MarketRunnerSearch(ctx, &request)
    
    for market := range markets {
        // Do something with market   
    }

    for err := range errCh {
        // Handle error   
    }
}
```
