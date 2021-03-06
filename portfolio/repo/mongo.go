package repo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mongoRepository struct {
	ctx        context.Context
	Client     *mongo.Client
	Collection *mongo.Collection
}

func NewMongoRepository(ctx context.Context) Repository {
	return &mongoRepository{
		ctx: ctx,
	}
}

func (m *mongoRepository) Connect() error {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return err
	}
	m.Client = client
	// define context
	// ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	// defer cancel()
	fmt.Println("Initializing Repository: Connecting to 'mongodb://localhost:27017'")
	err = m.Client.Connect(m.ctx)
	if err != nil {
		return err
	}

	m.Collection = m.Client.Database("portfoliodb").Collection("blog")
	return nil
}
func (m *mongoRepository) CreateBlog(req *Blog) (*Blog, error) {

	if req == nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Invalid argument: request is nil",
		)
	}

	res, err := m.Collection.InsertOne(m.ctx, *req)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal Error: %v", err),
		)
	}

	if returnedId, ok := res.InsertedID.(primitive.ObjectID); ok {
		return &Blog{
			ID:      returnedId,
			Content: req.Content,
			Title:   req.Title,
		}, nil
	}
	return nil, status.Errorf(
		codes.Internal,
		"Internal Error: Unable to parse the blog id",
	)
}

func (m *mongoRepository) ReadBlog(blogID string) (*Blog, error) {
	oid, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Cannot parse blog id")
	}
	response := &Blog{}
	filter := primitive.M{"_id": oid}
	if err := m.Collection.FindOne(m.ctx, filter).Decode(response); err != nil {
		return nil, status.Errorf(codes.NotFound, "Cannot Find blog with id = %v", blogID)
	}
	return response, nil
}

func (m *mongoRepository) Close() error {
	fmt.Println("Closing Repository: Disconnecting from 'mongodb://localhost:27017'")
	return m.Client.Disconnect(m.ctx)
}
