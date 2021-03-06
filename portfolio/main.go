package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	portfoliopb "github.com/go-learn/portfolio/pb"
	"github.com/go-learn/portfolio/repo"
	"github.com/go-learn/portfolio/server"
	"google.golang.org/grpc"
)

func main() {
	// If we crash the server we get the date time and the file name
	log.SetFlags(log.LstdFlags | log.Llongfile)
	fmt.Println("Portfolio Service starting")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
		return
	}
	defer func() {
		fmt.Println("Closing the listener")
		lis.Close()
	}()

	//Get mongo Repository
	repository := repo.NewMongoRepository(context.Background())
	err = repository.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		repository.Close()
	}()

	s := grpc.NewServer()
	portfoliopb.RegisterBlogServiceServer(s, &server.Server{Repository: repository})

	go func() {
		fmt.Println("Starting Sever...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to server: %v", err)
		}
	}()
	defer func() {
		fmt.Println("Stopping the server")
		s.Stop()
	}()

	// Wait for Ctrl C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch
}
