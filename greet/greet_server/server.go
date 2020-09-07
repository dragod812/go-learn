package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/go-learn/greet/greet_server/adapter"
	"github.com/go-learn/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/internal/status"
)

type server struct{}

//This is to demonstrate Unary gRPC API
func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	log.Printf("Greet function was invoked with request := %v", req)
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	if firstName == "" {
		//status holds the gRPC errors
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Received a request with no first name",
		)
	}
	result := "Hello! Mr/Ms. " + firstName + " " + lastName
	res := greetpb.GreetResponse{
		Result: result,
	}
	return &res, nil
}

func (*server) GreetManyTimes(request *greetpb.GreetRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	firstName := request.GetGreeting().GetFirstName()
	lastName := request.GetGreeting().GetLastName()
	for i := 0; i < 10; i++ {
		result := "Hello, " + strconv.Itoa(i) + "th time," + firstName + " " + lastName
		res := greetpb.GreetResponse{
			Result: result,
		}
		stream.Send(&res)
		time.Sleep(5 * time.Second)
	}
	return nil
}

func (*server) GreetGroup(stream greetpb.GreetService_GreetGroupServer) error {
	result := ""
	previousDelimiter := ""
	previousName := "Hello"
	for i := 1; ; i++ {
		msg, err := stream.Recv()
		switch err {
		case nil:
			result += previousDelimiter + previousName
			log.Printf("Received people %v", i)
			if previousDelimiter == "" {
				previousDelimiter = "! "
			} else {
				previousDelimiter = ", "
			}
			previousName = msg.GetGreeting().GetFirstName() + " " + msg.GetGreeting().GetLastName()
		case io.EOF:
			switch previousDelimiter {
			case "":
				previousName += " no One :-("
			case "! ":
			default:
				previousDelimiter = " and "
			}
			previousName += "."
			result += previousDelimiter + previousName
			response := greetpb.GreetResponse{
				Result: result,
			}
			return stream.SendAndClose(&response)
		default:
			log.Fatalf("Error while receiving response from client : %v", err)
			continue
		}
	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	for i := 1; ; i++ {
		msg, err := stream.Recv()
		switch err {
		case nil:
			log.Printf("Received Person %v\n", i)
			firstName := msg.GetGreeting().GetFirstName()
			lastName := msg.GetGreeting().GetLastName()
			result := "Hello, " + firstName + " " + lastName + "!"
			sendErr := stream.Send(adapter.AdaptStringToGreetResponse(result))
			if sendErr != nil {
				log.Fatalf("Error while sending response to client: %v", sendErr)
			}
		case io.EOF:
			return nil
		default:
			log.Fatalf("Error while receiving message from client: %v", err)
		}
	}
}

func main() {
	fmt.Println("First Grpc server test")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
