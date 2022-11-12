package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	calculatorpb "github.com/smith-golang/grpc-test01/calculator/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	fmt.Println("Hello from client")
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial("localhost:10010", opts...)
	if err != nil {
		log.Fatalf("could not connect :%v", err)
	}

	c := calculatorpb.NewCalculatorServiceClient(conn)
	doCalculate(c)
	doCalculateStreaming(c)
	doCalculatorClientStreaming(c)

	defer conn.Close()
}

func doCalculate(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.SumRequest{
		FirstNumber:  3,
		SecondNumber: 4,
	}
	//from proto service => server => client
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling %v", err)
	}
	log.Printf("Response from sum : %v", res.SumResult)
}

func doCalculateStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("start server streaming")
	req := &calculatorpb.SumOneRequest{
		FirstNumber: 5,
		LastNumber:  7,
	}
	resStream, err := c.SumStreaming(context.Background(), req)
	if err != nil {
		log.Fatalf("err while callling from server %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			//we reach the end of the streaming
			break
		}
		if err != nil {
			log.Fatalf("error while calling streaming %v", err)
		}
		log.Printf("Sum is : %v", msg.SumResult)
	}
}

func doCalculatorClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Client streaming is starting")
	requests := []*calculatorpb.SumManyRequest{
		&calculatorpb.SumManyRequest{
			FirstNumber: 12,
			LastNumber:  13,
		},
		&calculatorpb.SumManyRequest{
			FirstNumber: 4,
			LastNumber:  5,
		},
		&calculatorpb.SumManyRequest{
			FirstNumber: 6,
			LastNumber:  13,
		},
		&calculatorpb.SumManyRequest{
			FirstNumber: 12,
			LastNumber:  4,
		},
		&calculatorpb.SumManyRequest{
			FirstNumber: 7,
			LastNumber:  8,
		},
	}
	stream, err := c.SumCStreaming(context.Background())
	if err != nil {
		log.Fatalf("error while reading streaming %v", err)
	}
	for _, req := range requests {
		fmt.Printf("sending request %v \n", req)
		stream.Send(req)
		time.Sleep(time.Second * 2)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response %v", err)
	}
	fmt.Printf("client calculator result are %v", res)
}
