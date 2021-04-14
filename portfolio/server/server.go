package server

import (
	"context"
	"fmt"

	"github.com/go-learn/portfolio/adapter"
	portfoliopb "github.com/go-learn/portfolio/pb"
	"github.com/go-learn/portfolio/repo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	Repository repo.Repository
}

func (s *Server) CreateBlog(ctx context.Context, req *portfoliopb.CreateBlogRequest) (*portfoliopb.CreateBlogResponse, error) {
	fmt.Printf("Create Blog request received: %v", req)
	if req == nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Invalid argument: request is nil",
		)
	}
	blog, err := s.Repository.CreateBlog(adapter.AdaptToBlogRepositoryModel(req))
	if err != nil {
		return nil, err
	}
	return &portfoliopb.CreateBlogResponse{
		Blog: adapter.AdaptToProtoModel(blog),
	}, nil
}

func (s *Server) ReadBlog(ctx context.Context, req *portfoliopb.ReadBlogRequest) (*portfoliopb.ReadBlogResponse, error) {
	fmt.Printf("Read Blog request received: %v", req)
	if req == nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Invalid argument: request is nil",
		)
	}
	blog, err := s.Repository.ReadBlog(req.GetBlogId())
	if err != nil {
		return nil, err
	}
	return &portfoliopb.ReadBlogResponse{Blog: adapter.AdaptToProtoModel(blog)}, nil
}
