package main

import (
	"context"
	"fmt"
	"log"

	"github.com/smith-golang/grpc-test01/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	fmt.Println("Hello Client")
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial("0.0.0.0:10000", opts...)
	if err != nil {
		log.Fatalf("Could not connect : %v", err)
	}

	c := greetpb.NewGreetServiceClient(conn)

	doGreeting(c)
	greetAgain(c)
	doLogin(c)

	defer conn.Close()
}

func doGreeting(c greetpb.GreetServiceClient) {
	fmt.Println("Strating doGreeting GRPC service")
	req := &greetpb.GreetRequest{
		Greetreq: &greetpb.Greeting{
			FirstName: "Smith",
			LastName:  "Golang",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC %v", err)
	}
	log.Printf("Response from Greet : %v", res.Result)
}

func greetAgain(c greetpb.GreetServiceClient) {
	fmt.Printf("starting greetagain g_RPC Service")
	req := &greetpb.GreetRequest{
		Greetreq: &greetpb.Greeting{
			FirstName: "John",
			LastName:  "Doe",
		},
	}
	res, err := c.GreetAgain(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling greet RPC %v", err)
	}
	log.Printf("Response from Greet Again :%v", res.Result)
}

func doLogin(c greetpb.GreetServiceClient) {
	fmt.Println("Starting dologin g_RPC Service")
	req := &greetpb.LoginRequest{
		Loginreq: &greetpb.Logining{
			Username: "Smith",
			Password: "112233",
		},
	}
	res, err := c.Login(context.Background(), req)
	if err != nil {
		log.Fatalf("err while calling greet grpc %v", err)
	}
	log.Printf("Response from Greet : %v", res.Result)
}
