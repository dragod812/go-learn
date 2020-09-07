package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/go-learn/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("First grpc client test")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	// greetServer(c)
	// greetServerForStream(c)
	// greetMyGroup(c)
	greetEveryone(c)
}

func readName() (string, string) {
	fmt.Print("Name: ")
	firstName, lastName := "", ""
	fmt.Scanf("%v %v\n", &firstName, &lastName)
	return firstName, lastName
}

func getGreetRequest() *greetpb.GreetRequest {
	firstName, lastName := readName()
	req := greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: firstName,
			LastName:  lastName,
		},
	}
	return &req
}

func greetServer(c greetpb.GreetServiceClient) {
	fmt.Printf("Client created %f", c)
	req := getGreetRequest()
	resp, err := c.Greet(context.Background(), req)
	if err != nil {
		// ok is true if the error type is a gRPC defined error
		respErr, ok := status.FromError(err)
		if ok {
			if respErr.Code() == codes.InvalidArgument {
				log.Fatalf("Why you pass empty name maaan!")
			}
			return
		} else {
			log.Fatalf("Error while calling Greet RPC: %v", err)
			return
		}
	}
	log.Printf("Response from Greet: %v", resp.GetResult())
}

func greetServerForStream(c greetpb.GreetServiceClient) {
	fmt.Printf("Greet server for stream response")
	req := getGreetRequest()
	respStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes RPC: %v", err)
	}
	for {
		msg, err := respStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error While receiving stream: %v", err)
		}
		log.Printf("Stream message from server: %v", msg)
	}
}

func getGroupNames() []*greetpb.GreetRequest {
	var n int
	result := []*greetpb.GreetRequest{}
	fmt.Scanf("Number of People in group: %i\n", &n)
	for n > 0 {
		firstName, lastName := readName()
		greetRequest := greetpb.GreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: firstName,
				LastName:  lastName,
			},
		}
		result = append(result, &greetRequest)
		n--
	}
	return result
}

func greetMyGroup(c greetpb.GreetServiceClient) {

	stream, err := c.GreetGroup(context.Background())
	if err != nil {
		log.Fatalf("Error while calling GreetGroup: %v", err)
	}

	var n int
	fmt.Print("Number of People in group: ")
	fmt.Scanf("%d\n", &n)
	for n > 0 {
		firstName, lastName := readName()
		greetRequest := &greetpb.GreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: firstName,
				LastName:  lastName,
			},
		}
		fmt.Printf("Sending Request: %v\n", greetRequest)
		err := stream.Send(greetRequest)
		if err != nil {
			log.Fatalf("Error while sending request %v Error: %v", greetRequest, err)
		}
		n--
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from GreetGroup: %v\n", err)
	}
	fmt.Println(response)
}

func greetEveryone(c greetpb.GreetServiceClient) {

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while calling GreetEveryone: %v", err)
	}

	//make go routines to print responses from backend
	go func() {
		for {
			res, err := stream.Recv()
			switch err {
			case nil:
				fmt.Println(res)
			case io.EOF:
				return
			default:
				log.Fatalf("Error while receiving message from server: %v", err)
			}
		}
	}()

	for {
		firstName, lastName := readName()
		if firstName == "Done" {
			break
		}
		greetRequest := &greetpb.GreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: firstName,
				LastName:  lastName,
			},
		}
		fmt.Printf("Sending Request: %v\n", greetRequest)
		err := stream.Send(greetRequest)
		time.Sleep(1 * time.Second)
		if err != nil {
			log.Fatalf("Error while sending request %v Error: %v", greetRequest, err)
		}
	}
}
