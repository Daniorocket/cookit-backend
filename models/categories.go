package models

import (
	"context"
	"fmt"
	"time"

	"github.com/Daniorocket/cookit-backend/lib"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collectionCategory = "categories"

type Category struct {
	ID   string   `json:"id" bson:"id"`
	Name string   `json:"name" bson:"name" validate:"nonnil,nonzero"`
	File lib.File `json:"file" bson:"file" validate:"nonnil,nonzero"`
}

type MongoCategoryRepository struct {
	DbPointer    *mongo.Client
	DatabaseName string
}

type PostgreSQLCategoryRepository struct {
	ConStr       string
	DatabaseName string
}

func (d *MongoCategoryRepository) GetByID(id string) (Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := d.DbPointer.Database(d.DatabaseName).Collection("categories")
	cat := Category{}
	if err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&cat); err != nil {
		return Category{}, err
	}
	return cat, nil
}

func (d *MongoCategoryRepository) Create(c Category) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := d.DbPointer.Database(d.DatabaseName).Collection(collectionCategory)
	if _, err := collection.InsertOne(ctx, &c); err != nil {
		return err
	}
	return nil
}
func (d *MongoCategoryRepository) GetAll(page, limit int) ([]Category, int64, error) {
	var category Category
	categories := []Category{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := d.DbPointer.Database(d.DatabaseName).Collection(collectionCategory)
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64((page - 1) * limit))
	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	for cursor.Next(ctx) {
		if err = cursor.Decode(&category); err != nil {
			return nil, 0, err
		}
		categories = append(categories, category)
	}
	totalElements, err := collection.CountDocuments(ctx, bson.M{}, nil)
	if err != nil {
		return nil, 0, err
	}
	return categories, totalElements, nil
}

func (p *PostgreSQLCategoryRepository) Create(c Category) error {
	fmt.Println("Not implemented yet")
	return nil
}

func (p *PostgreSQLCategoryRepository) GetAll(page, limit int) ([]Category, int64, error) {
	fmt.Println("Not implemented yet")
	return nil, 0, nil
}

func (p *PostgreSQLCategoryRepository) GetByID(id string) (Category, error) {
	fmt.Println("Not implemented yet")
	return Category{}, nil
}
