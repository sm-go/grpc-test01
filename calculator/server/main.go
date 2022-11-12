package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

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

// for server steaming
func (*server) SumStreaming(req *calculatorpb.SumOneRequest, stream calculatorpb.CalculatorService_SumStreamingServer) error {
	firstnumber := req.FirstNumber
	lastnumber := req.LastNumber
	count := 10
	for i := 0; i < count; i++ {
		result := firstnumber + lastnumber
		res := &calculatorpb.SumManyResponse{
			SumResult: result,
		}
		stream.Send(res)
		time.Sleep(time.Second * 1)
	}
	return nil
}

// for client streaming
func (*server) SumCStreaming(stream calculatorpb.CalculatorService_SumCStreamingServer) error {
	fmt.Println("client streaming is starting")
	var result int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			//we have finished the client streaming
			return stream.SendAndClose(&calculatorpb.SumOneResponse{
				SumResult: result,
			})
		}
		if err != nil {
			log.Fatalf("error while reading %v", err)
		}
		firstNumber := req.GetFirstNumber()
		lastNumber := req.GetLastNumber()
		result += firstNumber + lastNumber
	}
}

// for average streaming
func (*server) AverageService(stream calculatorpb.CalculatorService_AverageServiceServer) error {
	fmt.Println("average client streaming is starting")
	var count int32
	var sum int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			//
			average := float64(sum) / float64(count)
			return stream.SendAndClose(&calculatorpb.AverageResponse{
				Average: average,
			})
		}
		if err != nil {
			log.Fatalf("error while reading client streaming : %v", err)
		}
		sum += req.GetNumber()
		count++
	}
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
