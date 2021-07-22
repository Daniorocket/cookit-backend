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

type MongoAuthRepository struct {
	DbPointer    *mongo.Client
	DatabaseName string
}

func (m *MongoAuthRepository) Register(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := m.DbPointer.Database(m.DatabaseName).Collection(collectionUsers)
	_, err := collection.InsertOne(ctx, &user)
	if err != nil {
		return ErrCreateUser
	}
	return nil
}
func (m *MongoAuthRepository) GetPassword(login string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := m.DbPointer.Database(m.DatabaseName).Collection(collectionUsers)
	user := User{}
	if err := collection.FindOne(ctx, bson.M{"username": login}).Decode(&user); err == nil {
		return user.Password, nil
	}
	if err := collection.FindOne(ctx, bson.M{"email": login}).Decode(&user); err == nil {
		return user.Password, nil
	}
	return "", errors.New("Failed to find user with passed login")
}

func (m *MongoAuthRepository) GetUserinfo(username string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := m.DbPointer.Database(m.DatabaseName).Collection(collectionUsers)
	user := User{}
	if err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user); err != nil {
		return User{}, ErrFindUser
	}
	user.Password = "" //Encoded password can't be sent
	return user, nil
}
