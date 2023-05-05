package models

import (
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// User is the data model for documents from the 'users' collection
type User struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name" validate:"required,gt=2"`
	Email     string             `json:"email" bson:"email" validate:"required,email"`
	Password  []byte             `json:"-" bson:"password"`
	SPassword string             `json:"password,omitempty" bson:"-" validate:"required,gte=8"`
	CPassword string             `json:"confirm_password,omitempty" bson:"-"`
	Score     int                `json:"score" bson:"-"`
	CreatedAt time.Time          `json:"cat" bson:"cat"`
	UpdatedAt time.Time          `json:"uat" bson:"uat"`
}

// nameUnique checks to see if there is already a user in the database with that name
func (x *User) nameUnique() bool {
	filter := bson.M{
		"name": x.Name,
	}
	err := users.FindOne(ctx, filter).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return true
		}
	}
	return false
}

// emailUnique checks to see if there is already a user in the database with that email
func (x *User) emailUnique() bool {
	filter := bson.M{
		"email": x.Email,
	}

	log.Print("finding user by email")
	err := users.FindOne(ctx, filter).Err()
	if err != nil {
		log.Print("error")
		if err == mongo.ErrNoDocuments {
			log.Print("err no document")
			return true
		}
		log.Println("!!--!!  error trying to check the email's uniqueness", err.Error())
	}
	return false
}

// Save a user to the database
func (x *User) Save() error {

	if err := validate.Struct(x); err != nil {
		return err
	}
	if !x.nameUnique() && !x.emailUnique() {
		return errors.New("name or email already taken")
	}
	x.CreatedAt = time.Now()
	x.UpdatedAt = time.Now()

	hash, err := x.generatePasswordHash()
	if err != nil {
		return err
	}
	x.Password = hash

	res, err := users.InsertOne(ctx, x)
	if err != nil {
		return err
	}

	x.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

// GetAll returns a list of all the users in the database
func (x User) GetAll() ([]*User, error) {
	var a []*User
	cur, err := users.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var b *User
		if err := cur.Decode(&b); err != nil {
			return nil, err
		}
		a = append(a, b)
	}

	return a, nil
}

// Populate takes a user with only an id and uses it to fill out all other fields
func (x User) Get() (*User, error) {

	err := users.FindOne(ctx, bson.M{"_id": x.ID}).Decode(&x)
	return &x, err
}

// Update a user
func (x *User) Update() error {

	filter := bson.M{
		"_id": x.ID,
	}
	update := bson.M{
		"$set": bson.M{
			"name":      x.Name,
			"updatedAt": time.Now(),
		},
	}

	err := users.FindOneAndUpdate(ctx, filter, update).Decode(&x)
	return err
}

// Destroy a user
func (x *User) Destroy() error {

	filter := bson.M{
		"_id": x.ID,
	}
	err := users.FindOneAndDelete(ctx, filter).Decode(&x)
	return err
}

func (x User) DestroyAll() error {
	_, err := users.DeleteMany(ctx, bson.M{})
	return err
}

// Login a user
func (x *User) Login() error {
	filter := bson.M{
		"name": x.Name,
	}
	pw := x.SPassword
	if err := users.FindOne(ctx, filter).Decode(&x); err != nil {
		return err
	}
	if err := x.checkPasswordHash(pw); err != nil {
		return err
	}
	return nil
}

// generatePasswordHash generates a bcrypt hash of the given password
func (x *User) generatePasswordHash() ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(x.SPassword), bcrypt.DefaultCost)
}

// checkPasswordHash checks if the given password matches the bcrypt hash
func (x *User) checkPasswordHash(password string) error {
	return bcrypt.CompareHashAndPassword(x.Password, []byte(password))
}
