package main

import (
	"context"
	"io"
	"log"

	calculatorPB "besg-grpc/proto/calculator"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	defer conn.Close()

	client := calculatorPB.NewCalculatorServiceClient(conn)
	sendRequest(client)
	doServerStreaming(6, client)
}

func sendRequest(client calculatorPB.CalculatorServiceClient) {
	req := &calculatorPB.CalculatorRequest{
		A: 3,
		B: 10,
	}
	res, err := client.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling CalculatorService: %v \n", err)
	}

	log.Printf("Response from CalculatorService: %v", res.Result)
}

func doServerStreaming(target int, client calculatorPB.CalculatorServiceClient) {
	req := &calculatorPB.GetFibonacciRequest{
		Num: 5,
	}
	stream, err := client.GetFibonacci(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling CalculatorService: %v \n", err)
	}

	for {
		msg, err := stream.Recv()

		// 表示 server 的 stream 資料傳完了
		if err == io.EOF {
			break
		}

		// 表示有錯誤發生
		if err != nil {
			log.Fatalf("error while receiving sever stream: %v", err)
		}

		// 讀取 server stream 傳來的資料
		log.Printf("Response from GetFibonacci: %v", msg.GetNum())
	}
}
