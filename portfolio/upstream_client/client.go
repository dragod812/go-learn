package main

import (
	"context"
	"fmt"
	"log"

	portfoliopb "github.com/go-learn/portfolio/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Portfolio Client examples")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := portfoliopb.NewBlogServiceClient(cc)
	// createBlog("Blog1", "BlogContent", c)
	// greetServer(c)
	// greetServerForStream(c)
	// greetMyGroup(c)
	// greetEveryone(c)
	readBlogRequest := portfoliopb.ReadBlogRequest{BlogId: "603b940aa3b257cce609d9ed"}
	response, err := c.ReadBlog(context.Background(), &readBlogRequest)
	fmt.Printf("Read blog id: %v \n Response: %v", readBlogRequest.GetBlogId(), response)
}

func createBlog(title string, content string, c portfoliopb.BlogServiceClient) {
	req := portfoliopb.CreateBlogRequest{
		Blog: &portfoliopb.Blog{Title: title, Content: content},
	}
	resp, err := c.CreateBlog(context.Background(), &req)
	if err != nil {
		// ok is true if the error type is a gRPC defined error
		respErr, ok := status.FromError(err)
		if ok {
			if respErr.Code() == codes.InvalidArgument {
				log.Fatalf("Empty Blog Details")
			}
			return
		} else {
			log.Fatalf("Error while calling Greet RPC: %v", err)
			return
		}
	}
	log.Printf("Response from Portfolio Service: %v", resp)
}
