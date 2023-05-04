package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Parent    primitive.ObjectID `json:"_parent" bson:"_parent"`
	Author    primitive.ObjectID `json:"_author" bson:"_author"`
	Body      string             `json:"body" bson:"body"`
	Score     int                `json:"score" bson:"-"`
	CreatedAt time.Time          `json:"createdAt" bson:"updatedAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

func (x *Comment) SetScore(v bool) {
	if v {
		x.Score += 1
	} else {
		x.Score -= 1
	}
}

func (x *Comment) GetID() primitive.ObjectID {
	return x.ID
}

func (x *Comment) Save() error {
	x.CreatedAt = time.Now()
	x.UpdatedAt = time.Now()

	res, err := comments.InsertOne(ctx, x)
	if err != nil {
		return err
	}
	x.ID = res.InsertedID.(primitive.ObjectID)

	return nil
}

func (x Comment) GetAll() ([]*Comment, error) {
	var a []*Comment
	cur, err := comments.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	for cur.Next(ctx) {
		var b *Comment
		if err := cur.Decode(&b); err != nil {
			return nil, err
		}

		calculateScore(b)
		a = append(a, b)
	}
	return a, nil
}

func (x Comment) Populate() (*Comment, error) {
	err := comments.FindOne(ctx, bson.M{"_id": x.ID}).Decode(&x)

	calculateScore(&x)
	return &x, err
}

func (x *Comment) Update() error {
	filter := bson.M{
		"id": x.ID,
	}
	update := bson.M{
		"body":      x.Body,
		"updatedAt": time.Now(),
	}
	err := comments.FindOneAndUpdate(ctx, filter, update).Err()
	return err
}

func (x *Comment) Destroy() error {

	filter := bson.M{
		"id": x.ID,
	}
	err := comments.FindOneAndDelete(ctx, filter).Err()
	return err
}
