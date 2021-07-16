package models

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/Daniorocket/cookit-back/lib"
	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var CollectionCategory = "categories"
var ErrCreateCategory = errors.New("Failed to create category record")
var ErrFindCategory = errors.New("Failed to find category record")

type Category struct {
	ID      string   `json:"id" bson:"id"`
	LabelPL string   `json:"labelPL" bson:"label_pl"`
	LabelEN string   `json:"labelEN" bson:"label_en"`
	File    lib.File `json:"file" bson:"file"`
}

func GetAllCategories(client *mongo.Client, db, page, limit string) ([]Category, int64, error) {
	var category Category
	categories := []Category{}
	p, err := strconv.Atoi(page)
	if err != nil {
		return nil, 0, err
	}
	l, err := strconv.Atoi(limit)
	if err != nil {
		return nil, 0, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := client.Database(db).Collection(CollectionCategory)
	findOptions := options.Find()
	findOptions.SetLimit(int64(l))
	findOptions.SetSkip(int64((p - 1) * l))
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
func CreateCategory(client *mongo.Client, db string, category Category) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := client.Database(db).Collection(CollectionCategory)
	if _, err := collection.InsertOne(ctx, &category); err != nil {
		return ErrCreateCategory
	}
	return nil
}
func GetCategoryByID(client *mongo.Client, db, id string) (Category, error) {
	var category Category
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := client.Database(db).Collection(CollectionCategory)
	cursor := collection.FindOne(ctx, bson.M{"id": id})
	if err := cursor.Decode(&category); err != nil {
		return Category{}, err
	}
	return category, nil
}
