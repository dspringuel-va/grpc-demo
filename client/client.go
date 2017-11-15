package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"

	fibonacci "github.com/dspringuel-va/grpc-demo/protos"

	"google.golang.org/grpc"
)

var (
	n   = flag.Int("n", 0, "Wanted Fibonacci Number")
	all = flag.Bool("a", false, "Is streaming all numbers")
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial("localhost:4678", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Fail to dial: %v", err)
	}
	defer conn.Close()

	client := fibonacci.NewFibonnaciServiceClient(conn)

	if !*all {
		fibonacciResponse, fibErr := client.GetFibonnaciNumber(context.Background(), &fibonacci.FibonacciRequest{N: int32(*n)})
		fmt.Printf("GetFibonacciNumber(%d): %d (%v)", *n, fibonacciResponse.FN, fibErr)

	} else {
		fibonacciResponseStream, fibErr := client.GetAllFibonacciNumbers(context.Background(), &fibonacci.FibonacciRequest{N: int32(*n)})
		fmt.Printf("GetAllFibonacciNumbers(%d) %v: ", *n, fibErr)
		for {
			fibResponse, err := fibonacciResponseStream.Recv()
			if err == io.EOF {
				break
			}
			fmt.Printf("%d ", fibResponse.FN)
		}
	}
}
