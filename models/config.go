package models

import (
	"context"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ctx      context.Context
	users    *mongo.Collection
	posts    *mongo.Collection
	comments *mongo.Collection
	votes    *mongo.Collection
	validate = validator.New()
)

func Init(x *mongo.Client) {
	ctx = context.Background()
	users = x.Database("thoughtsea").Collection("users")
	posts = x.Database("thoughtsea").Collection("posts")
	comments = x.Database("thoughtsea").Collection("comments")
	votes = x.Database("thoughtsea").Collection("votes")
}

type Vote struct {
	Voter primitive.ObjectID `json:"voter" bson:"voter"`
	Val   int8               `json:"val" bson:"val"`
}
