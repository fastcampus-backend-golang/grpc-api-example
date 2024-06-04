package main

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/madeindra/stock-grpc/data"
	pb "github.com/madeindra/stock-grpc/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type stockService struct {
	pb.UnimplementedStockServiceServer
}

func (s *stockService) ListStocks(ctx context.Context, _ *emptypb.Empty) (*pb.StockCodes, error) {
	configs := data.GetStockConfig()

	// ambil semua kode saham
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

	// aktifkan / nonaktifkan saham sesuai request
	for code, isEnabled := range toggles {
		data.ToggleStock(code, isEnabled)
	}

	// ambil semua kode saham
	configs := data.GetStockConfig()

	// filter hanya yang aktif
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
	// ambil semua kode saham
	configs := data.GetStockConfig()

	// filter hanya yang aktif
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
	// panggil di goroutine untuk memproses request selama koneksi masih terbuka
	go func(stream pb.StockService_LiveStockServer) {
		for {
			req, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				return
			}

			code := req.StockCode
			isEnabled := req.IsEnabled

			data.ToggleStock(code, isEnabled)
		}
	}(stream)

	// loop untuk mengirimkan data harga saham setiap detik
	for {
		select {
		case <-stream.Context().Done():
			return nil
		default:
			// pause setiap detik untuk simulasi harga saham yang berubah
			time.Sleep(1 * time.Second)

			configs := data.GetStockConfig()

			for code, isEnabled := range configs {
				if isEnabled {
					history := data.GetStockPrice(code)
					if len(history) == 0 {
						continue
					}

					latestPrice := history[len(history)-1]

					stream.Send(&pb.StockPrices{
						StockPrices: map[string]*pb.StockPrice{
							code: {
								Price:     latestPrice.Price,
								Timestamp: timestamppb.New(latestPrice.Timestamp),
							},
						},
					})
				}
			}
		}
	}
}
