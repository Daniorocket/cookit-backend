package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Daniorocket/cookit-backend/lib"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var collectionUsers = "users"
var ErrCreateUser = errors.New("failed to create user record")
var ErrFindUser = errors.New("failed to find user record")

type User struct {
	ID               string   `json:"id" bson:"id"`
	Username         string   `json:"username" bson:"username"`
	Password         string   `json:"-" bson:"password"`
	AvatarURL        string   `json:"avatarURL" bson:"avatar_url"`
	Email            string   `json:"email" bson:"email"`
	Description      string   `json:"description" bson:"description"`
	PasswordRemindID string   `json:"-" bson:"password_remind_id"`
	FavoritesRecipes []Recipe `json:"favoritesRecipes" bson:"favorites_recipes"`
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
type Email struct {
	Email string `json:"email" bson:"email" validate:"nonnil,nonzero"`
}
type Password struct {
	Password string `json:"password" bson:"password" validate:"nonnil,nonzero"`
}

type MongoAuthRepository struct {
	DbPointer    *mongo.Client
	DatabaseName string
}

type PostgreSQLAuthRepository struct {
	ConStr       string
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
	return "", errors.New("failed to find user with passed login")
}

func (m *MongoAuthRepository) GetUserinfo(username string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := m.DbPointer.Database(m.DatabaseName).Collection(collectionUsers)
	user := User{}
	if err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user); err != nil {
		return User{}, err
	}
	return user, nil
}
func (m *MongoAuthRepository) CheckEmail(email string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := m.DbPointer.Database(m.DatabaseName).Collection(collectionUsers)
	user := User{}
	if err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return User{}, err
	}
	return user, nil
}
func (m *MongoAuthRepository) Update(userID string, user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := m.DbPointer.Database(m.DatabaseName).Collection(collectionUsers)
	res, err := collection.UpdateOne(ctx, bson.M{"id": userID},
		bson.M{"$set": bson.M{
			"email":              user.Email,
			"password":           user.Password,
			"avatar_url":         user.AvatarURL,
			"description":        user.Description,
			"password_remind_id": user.PasswordRemindID,
			"username":           user.Username,
			"favorites_recipes":  user.FavoritesRecipes,
		}})
	if err != nil {
		return err
	}
	fmt.Printf("Updated %v Documents!\n", res.ModifiedCount)
	return nil
}
func (m *MongoAuthRepository) GetUserByPasswordRemindID(passwordRemindID string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := m.DbPointer.Database(m.DatabaseName).Collection(collectionUsers)
	user := User{}
	if err := collection.FindOne(ctx, bson.M{"password_remind_id": passwordRemindID}).Decode(&user); err != nil {
		return User{}, err
	}
	return user, nil
}
func (d *MongoAuthRepository) DeleteUserAccount(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := d.DbPointer.Database(d.DatabaseName).Collection(collectionUsers)
	result, err := collection.DeleteOne(ctx, bson.M{"id": userID})
	if err != nil {
		return err
	}
	fmt.Printf("deleted  %v document(s)\n", result.DeletedCount)
	return nil
}

func (p *PostgreSQLAuthRepository) CheckEmail(email string) (User, error) {
	fmt.Println("Not implemented yet")
	return User{}, nil
}

func (p *PostgreSQLAuthRepository) DeleteUserAccount(userID string) error {
	fmt.Println("Not implemented yet")
	return nil
}

func (p *PostgreSQLAuthRepository) GetPassword(login string) (string, error) {
	fmt.Println("Not implemented yet")
	return "", nil
}

func (p *PostgreSQLAuthRepository) GetUserByPasswordRemindID(passwordRemindID string) (User, error) {
	fmt.Println("Not implemented yet")
	return User{}, nil
}

func (p *PostgreSQLAuthRepository) GetUserinfo(username string) (User, error) {
	fmt.Println("Not implemented yet")
	return User{}, nil
}

func (p *PostgreSQLAuthRepository) Register(user User) error {
	fmt.Println("user:", user)
	db, err := lib.ConnectPostgres(p.ConStr)
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO users(id,username,password,email) VALUES($1,$2,$3,$4);")
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(user.ID, user.Username, user.Password, user.Email); err != nil {
		return err
	}
	return nil
}

func (p *PostgreSQLAuthRepository) Update(userID string, user User) error {
	fmt.Println("Not implemented yet")
	return nil
}
