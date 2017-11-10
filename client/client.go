package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	fibonacci "github.com/dspringuel-va/grpc-demo/protos"

	"google.golang.org/grpc"
)

var (
	n = flag.Int("n", 0, "Wanted Fibonacci Number")
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial("localhost:4678", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Fail to dial: %v", err)
	}
	defer conn.Close()

	client := fibonacci.NewFibonnaciServiceClient(conn)

	fibonacciResponse, fibErr := client.GetFibonnaciNumber(context.Background(), &fibonacci.FibonacciRequest{N: int32(*n)})

	fmt.Printf("GetFibonacciNumber(%d): %d (%v)", *n, fibonacciResponse.FN, fibErr)
}
