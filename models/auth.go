package models

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var collectionUsers = "users"
var ErrCreateUser = errors.New("Failed to create user record")
var ErrFindUser = errors.New("Failed to find user record")

type User struct {
	ID          string `json:"id" bson:"id"`
	Username    string `json:"username" bson:"username"`
	Password    string `json:"password" bson:"password"`
	AvatarURL   string `json:"avatarURL" bson:"avatar_url"`
	Email       string `json:"email" bson:"email"`
	Description string `json:"description" bson:"description"`
}
type Credentials struct {
	Email    string `json:"email" bson:"email"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}
type Login struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

func RegisterUser(client *mongo.Client, db string, u *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := client.Database(db).Collection(collectionUsers)
	_, err := collection.InsertOne(ctx, &u)
	if err != nil {
		return ErrCreateUser
	}
	return nil
}
func GetPasswordByUsernameOrEmail(client *mongo.Client, db string, username string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := client.Database(db).Collection(collectionUsers)
	user := User{}
	if err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user); err == nil {
		return user.Password, nil
	}
	if err := collection.FindOne(ctx, bson.M{"email": username}).Decode(&user); err == nil {
		return user.Password, nil
	}
	return "", ErrFindUser
}
func GetUserinfo(client *mongo.Client, db string, username string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := client.Database(db).Collection(collectionUsers)
	user := User{}
	if err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user); err != nil {
		return User{}, ErrFindUser
	}
	user.Password = "" //Encoded password can't be sended
	return user, nil
}
