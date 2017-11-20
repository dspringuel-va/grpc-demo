package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"

	fibonacci "github.com/dspringuel-va/grpc-demo/protos"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type fibonacciServer struct {
}

func (fibServer *fibonacciServer) GetFibonnaciNumber(ctx context.Context, request *fibonacci.FibonacciRequest) (*fibonacci.FibonacciResponse, error) {
	fmt.Printf("\nGetFibonnaciNumber: %d\n", request.N)
	return &fibonacci.FibonacciResponse{FN: fibNumber(request.N)}, nil
}

func (fibServer *fibonacciServer) GetAllFibonacciNumbers(request *fibonacci.FibonacciRequest, stream fibonacci.FibonnaciService_GetAllFibonacciNumbersServer) error {
	fmt.Printf("\nGetAllFibonacciNumbers: %d\n", request.N)
	for i := int32(0); i <= request.N; i++ {
		fmt.Printf("Sending FN(%d)\n", i)
		if err := stream.Send(&fibonacci.FibonacciResponse{FN: fibNumber(i)}); err != nil {
			return err
		}
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

func (fibServer *fibonacciServer) JoinFibonacciNumbers(stream fibonacci.FibonnaciService_JoinFibonacciNumbersServer) error {
	fmt.Printf("\nJoinFibonacciNumbers\n")
	var clientFibNumbers []string

	for {
		fibRequest, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("\nClient has stopped streaming\n")
			return stream.SendAndClose(&fibonacci.JoinedFibonacciResponse{JoinedFN: strings.Join(clientFibNumbers, " - ")})
		}

		if err != nil {
			return err
		}

		fmt.Printf("Joining %d\n", fibNumber(fibRequest.N))
		clientFibNumbers = append(clientFibNumbers, fmt.Sprintf("%d", fibNumber(fibRequest.N)))
	}
}

func (fibServer *fibonacciServer) ElevatorFibonacci(stream fibonacci.FibonnaciService_ElevatorFibonacciServer) error {
	fmt.Printf("\nElevatorFibonacci\n")

	i := int32(10)
	mod := int32(1)
	go func(stream fibonacci.FibonnaciService_ElevatorFibonacciServer) {
		for {
			if newI := i + mod; newI >= 0 && newI <= 20 {
				i = newI
				fmt.Printf("Sending FN(%d)\n", i)
				if err := stream.Send(&fibonacci.FibonacciResponse{FN: fibNumber(i)}); err != nil {
					fmt.Printf("Client has stopped streaming (%v)\n", err)
					return
				}
			}

			time.Sleep(time.Duration(50+rand.Intn(300)) * time.Millisecond)
		}
	}(stream)

	for {
		_, err := stream.Recv()

		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		fmt.Printf("Receiving Switch\n")
		mod = mod * -1
	}
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
