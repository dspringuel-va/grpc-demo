package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	fibonacci "github.com/dspringuel-va/grpc-demo/protos"

	"google.golang.org/grpc"
)

var (
	n    = flag.Int("n", 0, "Wanted Fibonacci Number")
	all  = flag.Bool("a", false, "Is streaming all numbers")
	join = flag.Bool("j", false, "Is server joined client streamed number")
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial("localhost:4678", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Fail to dial: %v", err)
	}
	defer conn.Close()

	client := fibonacci.NewFibonnaciServiceClient(conn)

	if *all {
		fibonacciResponseStream, fibErr := client.GetAllFibonacciNumbers(context.Background(), &fibonacci.FibonacciRequest{N: int32(*n)})
		fmt.Printf("GetAllFibonacciNumbers(%d) %v: ", *n, fibErr)
		for {
			fibResponse, err := fibonacciResponseStream.Recv()
			if err == io.EOF {
				break
			}
			fmt.Printf("%d ", fibResponse.FN)
		}
	} else if *join {
		fibonacciStream, fibErr := client.JoinFibonacciNumbers(context.Background())
		if fibErr != nil {
			log.Fatalf("Can't open client stream: %v", fibErr)
		}

		clientFibNumbers := [10]int32{4, 8, 3, 1, 7, 23, 16, 6, 5, 12}
		fmt.Printf("Sending")
		for _, n := range clientFibNumbers {
			fmt.Printf(" %d", n)
			fibonacciStream.Send(&fibonacci.FibonacciRequest{N: n})
			time.Sleep(500 * time.Millisecond)
		}

		joinedFib, joinErr := fibonacciStream.CloseAndRecv()

		fmt.Printf("\nJoinFibonacciNumbers: %s (%v)", joinedFib.JoinedFN, joinErr)

	} else {
		fibonacciResponse, fibErr := client.GetFibonnaciNumber(context.Background(), &fibonacci.FibonacciRequest{N: int32(*n)})
		fmt.Printf("GetFibonacciNumber(%d): %d (%v)", *n, fibonacciResponse.FN, fibErr)
	}
}
