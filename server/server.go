package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	calculatorPB "besg-grpc/proto/calculator"

	"google.golang.org/grpc"
)

// 實作方式需參考
// calculator_grpc.pb.go
type Server struct {
	calculatorPB.UnimplementedCalculatorServiceServer
}

// 存在 calculator_grpc.pb.go中
// type CalculatorServiceServer interface {
// 	Sum(context.Context, *CalculatorRequest) (*CalculatorResponse, error)
// 	// 5 -> 1 1 2 3 5
// 	GetFibonacci(*GetFibonacciRequest, CalculatorService_GetFibonacciServer) error
// 	mustEmbedUnimplementedCalculatorServiceServer()
// }
func (*Server) Sum(ctx context.Context, req *calculatorPB.CalculatorRequest) (*calculatorPB.CalculatorResponse, error) {
	fmt.Printf("Sum function is invoked with %v \n", req)

	a := req.GetA() // req.A
	b := req.GetB() // req.B

	res := &calculatorPB.CalculatorResponse{
		// 大寫開頭
		Result: a + b,
	}

	return res, nil
}

func (*Server) GetFibonacci(req *calculatorPB.GetFibonacciRequest, stream calculatorPB.CalculatorService_GetFibonacciServer) error {
	fmt.Printf("GetFibonacci function is invoked with %v \n", req)

	position := req.GetNum()
	cache := make([]int64, position+1)
	result := fibMemo(position, cache)

	for _, num := range result {
		stream.Send(&calculatorPB.GetFibonacciResponse{
			Num: int64(num),
		})
		time.Sleep(1 * time.Second)
	}

	return nil
}

func fibMemo(position int64, cache []int64) []int64 {
	if cache[position] != 0 {
		return cache
	} else {
		if position <= 2 {
			cache[position] = 1
		} else {
			cache[position] = fibMemo(position-1, cache)[position-1] + fibMemo(position-2, cache)[position-2]
		}

		return cache
	}
}

func main() {
	fmt.Println("starting gRPC server...")

	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v \n", err)
	}

	grpcServer := grpc.NewServer()
	calculatorPB.RegisterCalculatorServiceServer(grpcServer, &Server{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v \n", err)
	}
}
