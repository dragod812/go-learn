package repo

import "go.mongodb.org/mongo-driver/bson/primitive"

type Blog struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Content string             `bson:"content"`
	Title   string             `bson:"title"`
}

type Repository interface {
	Connect() error
	CreateBlog(*Blog) (*Blog, error)
	ReadBlog(blogID string) (*Blog, error)
	Close() error
}
