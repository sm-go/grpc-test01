package main

import (
	"context"
	"fmt"
	"io"
	"log"

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
