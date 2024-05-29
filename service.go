package main

import (
	"context"
	"errors"
	"io"

	"github.com/madeindra/stock-grpc/data"
	pb "github.com/madeindra/stock-grpc/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type stockService struct {
	pb.UnimplementedStockServiceServer
}

func (s *stockService) ListStocks(ctx context.Context, _ *emptypb.Empty) (*pb.StockCodes, error) {
	configs := data.GetStockConfig()

	// get all stock codes
	codes := []string{}
	for code := range configs {
		codes = append(codes, code)
	}

	reponse := &pb.StockCodes{
		StockCodes: codes,
	}

	return reponse, nil
}

func (s *stockService) ToggleStocks(stream pb.StockService_ToggleStocksServer) error {
	toggles := make(map[string]bool)

	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}

		code := req.StockCode
		isEnabled := req.IsEnabled

		toggles[code] = isEnabled
	}

	// bulk toggle stock
	for code, isEnabled := range toggles {
		data.ToggleStock(code, isEnabled)
	}

	// get latest stock config
	configs := data.GetStockConfig()

	// filter to get only the enabled stock
	subscribed := []string{}
	for code, isEnabled := range configs {
		if isEnabled {
			subscribed = append(subscribed, code)
		}
	}

	return stream.SendAndClose(&pb.StockCodes{
		StockCodes: subscribed,
	})
}

func (s *stockService) ListSubscriptions(_ *emptypb.Empty, stream pb.StockService_ListSubscriptionsServer) error {
	// get latest stock config
	configs := data.GetStockConfig()

	// filter to find only the enabled stock
	for code, isEnabled := range configs {
		if isEnabled {
			stream.Send(&pb.StockCode{
				StockCode: code,
			})
		}
	}

	return nil
}

func (s *stockService) LiveStock(stream pb.StockService_LiveStockServer) error {
	return nil
}
