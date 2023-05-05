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
	Votes     []Vote             `json:"-" bson:"votes"`
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
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var b *Comment
		if err := cur.Decode(&b); err != nil {
			return nil, err
		}
		b.calculateScore()
		a = append(a, b)
	}
	return a, nil
}

func (x Comment) Get() (*Comment, error) {
	if err := comments.FindOne(ctx, bson.M{"_id": x.ID}).Decode(&x); err != nil {
		return nil, err
	}
	x.calculateScore()
	return &x, nil

}

func (x *Comment) Update() error {
	filter := bson.M{
		"id": x.ID,
	}
	update := bson.M{
		"$set": bson.M{
			"body":      x.Body,
			"updatedAt": time.Now(),
		},
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

func (x *Comment) Vote(vote *Vote) error {

	filter := bson.M{
		"_id": x.ID,
	}
	update := bson.M{
		"$set": bson.M{
			"votes": x.Votes,
		},
	}

	i, exists := x.alreadyVoted(vote)
	if exists {
		if x.Votes[i].Val == vote.Val {
			// remove vote
			x.Votes = append(x.Votes[:i], x.Votes[i+1:]...)
		} else {
			// update vote
			x.Votes[i] = *vote
		}
	} else {
		// create vote
		x.Votes = append(x.Votes, *vote)
	}
	if err := comments.FindOneAndUpdate(ctx, filter, update).Err(); err != nil {
		return err
	}
	return nil
}

func (x *Comment) alreadyVoted(vote *Vote) (int, bool) {
	for i, v := range x.Votes {
		if v.Voter == vote.Voter {
			return i, true
		}
	}
	return 0, false
}

func (x *Comment) calculateScore() {
	for _, v := range x.Votes {
		x.Score += int(v.Val)
	}
}
