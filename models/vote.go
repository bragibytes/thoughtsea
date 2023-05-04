package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Vote struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Parent    primitive.ObjectID `json:"parent" bson:"parent"`
	Voter     primitive.ObjectID `json:"voter" bson:"voter"`
	Up        bool               `json:"isUpvote" bson:"isUpvote"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

// create a vote
func (x *Vote) create() error {

	x.CreatedAt = time.Now()
	x.UpdatedAt = time.Now()
	res, err := votes.InsertOne(ctx, x)
	if err != nil {
		return err
	}
	x.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

// Cast a vote
func (x *Vote) Cast() error {
	var v *Vote

	filter := bson.M{
		"parent": x.Parent,
		"voter":  x.Voter,
	}
	err := votes.FindOne(ctx, filter).Decode(&v)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// create vote
			if err := x.create(); err != nil {
				return err
			}
		}
		return err
	} else {
		// vote exists, either update or destroy it
		if x.Up != v.Up {
			//update
			v.Up = x.Up
			if err := v.update(); err != nil {
				return err
			}
		} else {
			//destroy
			if err := v.destroy(); err != nil {
				return err
			}
		}
	}

	return nil
}

// update a vote
func (x *Vote) update() error {

	filter := bson.M{
		"_id": x.ID,
	}
	update := bson.M{
		"isUpvote":  x.Up,
		"updatedAt": x.UpdatedAt,
	}

	err := votes.FindOneAndUpdate(ctx, filter, update).Err()

	return err
}

// destroy a vote
func (x *Vote) destroy() error {

	err := votes.FindOneAndDelete(ctx, bson.M{
		"_id": x.ID,
	}).Err()

	return err
}
