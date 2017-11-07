package main

import (
	"fmt"
	"log"
	"net"

	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

type fibonacciServer struct {
}

func (fibServer *fibonacciServer) GetFibonnaciNumber(ctx context.Context, request *FibonacciRequest) (*FibonacciResponse, error) {
	if request.GetN() == 0 {
		return &FibonacciResponse{FN: 0}, nil
	}
	if request.GetN() == 1 {
		return &FibonacciResponse{FN: 1}, nil
	}
	var fn int32 = 1
	var fnMinusOne int32
	var fnMinusTwo int32
	var i int32
	for i = 2; i <= request.N; i++ {
		fn, fnMinusOne, fnMinusTwo = fnMinusOne+fnMinusTwo, fn, fnMinusOne
	}

	return &FibonacciResponse{FN: fn}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:4678"))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grcpServer := grpc.NewServer()
	RegisterFibonnaciServiceServer(grcpServer, new(fibonacciServer))
	grcpServer.Serve(lis)
}
