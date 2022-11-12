package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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
	doServerStreaming(c)
	doClientStreaming(c)
	doBiDirectionalStreaming(c)

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

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting do server streaming from client")
	req := &greetpb.GreetOneRequest{
		Onereq: &greetpb.Greeting{
			FirstName: "Strange",
			LastName:  "Go",
		},
	}
	resStream, err := c.ServerGreeting(context.Background(), req)
	if err != nil {
		log.Fatalf("err while calling from server %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading streaming %v", err)
		}
		log.Printf("Response from %v", msg.GetResult())
	}
}

// for client streaming
func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting do client streaming from client")
	requests := []*greetpb.GreetManyRequest{
		&greetpb.GreetManyRequest{
			Manyreq: &greetpb.Greeting{
				FirstName: "Toe",
				LastName:  "Lin",
			},
		},
		&greetpb.GreetManyRequest{
			Manyreq: &greetpb.Greeting{
				FirstName: "Smith",
				LastName:  "Go",
			},
		},
		&greetpb.GreetManyRequest{
			Manyreq: &greetpb.Greeting{
				FirstName: "John",
				LastName:  "Doe",
			},
		},
		&greetpb.GreetManyRequest{
			Manyreq: &greetpb.Greeting{
				FirstName: "John one",
				LastName:  "Doe",
			},
		},
		&greetpb.GreetManyRequest{
			Manyreq: &greetpb.Greeting{
				FirstName: "John two",
				LastName:  "Doe",
			},
		},
	}

	stream, err := c.ClientGreet(context.Background())
	if err != nil {
		log.Fatalf("error while reading streaming %v", err)
	}
	for _, req := range requests {
		fmt.Printf("sending request %v", req)
		stream.Send(req)
		time.Sleep(time.Second * 2)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response %v", err)
	}
	fmt.Printf("Client greet %v \n", res)
}
func doBiDirectionalStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("bi-directional streaming from client")

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greetmany: &greetpb.Greeting{
				FirstName: "Toe",
				LastName:  "Lin",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greetmany: &greetpb.Greeting{
				FirstName: "Smith",
				LastName:  "Go",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greetmany: &greetpb.Greeting{
				FirstName: "John",
				LastName:  "Doe",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greetmany: &greetpb.Greeting{
				FirstName: "John one",
				LastName:  "Doe",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greetmany: &greetpb.Greeting{
				FirstName: "John two",
				LastName:  "Doe",
			},
		},
	}

	// create a stream by invoking the client
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream :%v", err)
		return
	}

	waitc := make(chan struct{})

	// send a bunch of messages to the client (go routine)
	go func() {
		for _, req := range requests {
			fmt.Printf("Sending message %v", req)
			stream.Send(req)
			time.Sleep(time.Second * 1)
		}
		stream.CloseSend()
	}()

	// receive a bunch of message from the client (go routine)
	go func() {
		res, err := stream.Recv()
		if err == io.EOF {
			close(waitc)
		}
		if err != nil {
			log.Fatalf("error while receiving %v", err)
			close(waitc)
		}
		fmt.Printf("receiving %v", res)
	}()

	// block until everything is done
	<-waitc
}
