package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	"google.golang.org/grpc/credentials"

	fibonacci "github.com/dspringuel-va/grpc-demo/protos"

	"google.golang.org/grpc"
)

var (
	n        = flag.Int("n", 0, "Wanted Fibonacci Number")
	all      = flag.Bool("a", false, "Is streaming all numbers")
	join     = flag.Bool("j", false, "Is server joined client streamed number")
	elevator = flag.Bool("e", false, "Elevator Fibonacci")
)

func main() {
	flag.Parse()

	creds, err := credentials.NewClientTLSFromFile("cert/domain.crt", "")
	if err != nil {
		log.Fatalf("Fail to create credentials: %v", err)
	}

	conn, err := grpc.Dial("localhost:4678", grpc.WithTransportCredentials(creds))

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

	} else if *elevator {

		stream, fibErr := client.ElevatorFibonacci(context.Background())
		fmt.Printf("Elevator (%v)\n", fibErr)

		waitCompleted := make(chan struct{})

		fmt.Printf("Response: \n")
		go func(stream fibonacci.FibonnaciService_ElevatorFibonacciClient) {
			for {
				fibResponse, err := stream.Recv()
				if err == io.EOF {
					fmt.Printf("\nServer ended streaming")
					close(waitCompleted)
					return
				}
				if err != nil {
					log.Fatalf("Elevator error %v", err)
				}
				fmt.Printf(" %d", fibResponse.FN)
			}
		}(stream)

		for i := int32(0); i < 5; i++ {
			time.Sleep(time.Duration(2000+rand.Intn(1000)) * time.Millisecond)
			fmt.Printf("\nSending switch\n")
			stream.Send(&fibonacci.SwitchRequest{})
		}

		time.Sleep(2 * time.Second)
		fmt.Printf("\n\nClose client connection")
		stream.CloseSend()
		<-waitCompleted
		fmt.Printf("\nElevator ended")

	} else {
		fibonacciResponse, fibErr := client.GetFibonnaciNumber(context.Background(), &fibonacci.FibonacciRequest{N: int32(*n)})
		if fibErr != nil {
			log.Fatalf("Error while getting fibonacci number: %v", fibErr)
		}
		fmt.Printf("GetFibonacciNumber(%d): %d", *n, fibonacciResponse.FN)
	}
}
