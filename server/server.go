package main

import (
	"fmt"
	"log"
	"net"
	"time"

	fibonacci "github.com/dspringuel-va/grpc-demo/protos"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type fibonacciServer struct {
}

func (fibServer *fibonacciServer) GetFibonnaciNumber(ctx context.Context, request *fibonacci.FibonacciRequest) (*fibonacci.FibonacciResponse, error) {
	return &fibonacci.FibonacciResponse{FN: fibNumber(request.N)}, nil
}

func (fibServer *fibonacciServer) GetAllFibonacciNumbers(request *fibonacci.FibonacciRequest, stream fibonacci.FibonnaciService_GetAllFibonacciNumbersServer) error {
	for i := int32(0); i <= request.N; i++ {
		if err := stream.Send(&fibonacci.FibonacciResponse{FN: fibNumber(i)}); err != nil {
			return err
		}
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

func fibNumber(n int32) int32 {
	if n == 0 || n == 1 {
		return n
	}

	fn := int32(1)
	fnMinusOne := int32(0)
	for i := int32(2); i <= n; i++ {
		fn, fnMinusOne = fn+fnMinusOne, fn
	}

	return fn
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:4678"))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grcpServer := grpc.NewServer()
	fibonacci.RegisterFibonnaciServiceServer(grcpServer, new(fibonacciServer))
	fmt.Println("Listening to port 4678")
	grcpServer.Serve(lis)
}
