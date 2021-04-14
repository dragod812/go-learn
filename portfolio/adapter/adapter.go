package adapter

import (
	portfoliopb "github.com/go-learn/portfolio/pb"
	"github.com/go-learn/portfolio/repo"
)

func AdaptToBlogRepositoryModel(req *portfoliopb.CreateBlogRequest) *repo.Blog {
	return &repo.Blog{
		Title:   req.GetBlog().GetTitle(),
		Content: req.GetBlog().GetContent(),
	}
}

func AdaptToProtoModel(req *repo.Blog) *portfoliopb.Blog {
	return &portfoliopb.Blog{
		Id:      req.ID.Hex(),
		Content: req.Content,
		Title:   req.Title,
	}
}
