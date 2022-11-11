package main

import (
	"context"
	"fmt"
	"log"
	"net"

	calculatorpb "github.com/smith-golang/grpc-test01/calculator/pb"
	"google.golang.org/grpc"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

// from proto service
func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	firstNumner := req.GetFirstNumber()
	secondNumber := req.GetSecondNumber()
	sum := firstNumner + secondNumber
	res := &calculatorpb.SumResponse{
		SumResult: sum,
	}
	return res, nil
}

func main() {
	fmt.Println("hello calculator server")
	list, err := net.Listen("tcp", "localhost:10010")
	if err != nil {
		log.Fatalf("failed to listen server %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})
	if err := s.Serve(list); err != nil {
		log.Fatalf("Failed to served %v", err)
	}
}
