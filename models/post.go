package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Author    primitive.ObjectID `json:"_author" bson:"_author"`
	Title     string             `json:"title" bson:"title"`
	Body      string             `json:"body" bson:"body"`
	Score     int                `json:"score" bson:"-"`
	Votes     []Vote             `json:"-" bson:"votes"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

func (x *Post) SetScore(v bool) {
	if v {
		x.Score += 1
	} else {
		x.Score -= 1
	}
}
func (x *Post) GetID() primitive.ObjectID {
	return x.ID
}

func (x *Post) Save() error {

	x.CreatedAt = time.Now()
	x.UpdatedAt = time.Now()

	if _, err := posts.InsertOne(ctx, x); err != nil {
		return err
	}
	return nil
}

func (x Post) GetAll() ([]*Post, error) {

	var a []*Post
	cur, err := posts.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var b *Post
		if err := cur.Decode(&b); err != nil {
			return nil, err
		}
		b.calculateScore()
		a = append(a, b)
	}
	return a, nil
}
func (x Post) Get() (*Post, error) {
	if err := posts.FindOne(ctx, bson.M{"_id": x.ID}).Decode(&x); err != nil {
		return nil, err
	}
	x.calculateScore()
	return &x, nil
}

func (x *Post) Update() error {

	filter := bson.M{
		"_id": x.ID,
	}
	update := bson.M{
		"$set": bson.M{
			"title":      x.Title,
			"body":       x.Body,
			"_updatedAt": time.Now(),
		},
	}

	err := posts.FindOneAndUpdate(ctx, filter, update).Decode(&x)

	return err
}

func (x *Post) Destroy() error {

	filter := bson.M{
		"_id": x.ID,
	}
	err := posts.FindOneAndDelete(ctx, filter).Decode(&x)
	return err
}

func (x *Post) Vote(vote *Vote) error {

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

	if err := posts.FindOneAndUpdate(ctx, filter, update).Err(); err != nil {
		return err
	}
	return nil
}

func (x *Post) alreadyVoted(vote *Vote) (int, bool) {
	for i, v := range x.Votes {
		if v.Voter == vote.Voter {
			return i, true
		}
	}
	return 0, false
}

func (x *Post) calculateScore() {
	for _, v := range x.Votes {
		x.Score += int(v.Val)
	}
}
