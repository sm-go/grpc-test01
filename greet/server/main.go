package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/smith-golang/grpc-test01/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("greeting functions was invoked with %v", req)
	firstName := req.GetGreetreq().GetFirstName()
	lastName := req.GetGreetreq().GetLastName()
	result := "hello " + firstName + lastName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func (*server) GreetAgain(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("greet again function was invoked %v", req)
	fName := req.GetGreetreq().FirstName
	lName := req.GetGreetreq().LastName
	result := fName + lName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func (*server) Login(ctx context.Context, req *greetpb.LoginRequest) (*greetpb.LoginResponse, error) {
	fmt.Printf("login functions was invoked with %v", req)
	username := req.GetLoginreq().GetUsername()
	password := req.GetLoginreq().GetPassword()
	result := "Login data :" + username + password
	res := &greetpb.LoginResponse{
		Result: result,
	}
	return res, nil
}

func main() {
	fmt.Println("Hello Server")

	listen, err := net.Listen("tcp", "0.0.0.0:10000")
	if err != nil {
		log.Fatalf("Failed to listen :%v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to served :%v", err)
	}

}
