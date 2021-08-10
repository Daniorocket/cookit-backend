package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collectionRecipes = "recipes"

type KitchenStyle int

const (
	Polish KitchenStyle = iota
	Russian
)

type Tags int

const (
	Easy Tags = iota
	Medium
	Hard
)

type TagsList struct {
	Tags []int `json:"tags"`
}

type Recipe struct {
	ID               string       `json:"id" bson:"id"`
	Name             string       `json:"name" bson:"name" validate:"nonnil,nonzero"`
	Username         string       `json:"username" bson:"username" validate:"nonnil,nonzero"`
	Kitchen          KitchenStyle `json:"kitchenStyle" bson:"kitchen_style" validate:"nonnil,nonzero"`
	Tags             []Tags       `json:"tags" bson:"tags" validate:"nonnil,nonzero"`
	ListOfSteps      []string     `json:"listOfSteps" bson:"list_of_steps" validate:"nonnil,nonzero"`
	ListOfCategories []Category   `json:"listOfCategories" bson:"list_of_categories" validate:"nonnil,nonzero"`
	Description      string       `json:"description" bson:"description" validate:"nonnil,nonzero"`
	Date             string       `json:"date" bson:"date" validate:"nonnil,nonzero"`
}

type MongoRecipeRepository struct {
	DbPointer    *mongo.Client
	DatabaseName string
}

func (m *MongoRecipeRepository) Create(recipe Recipe) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := m.DbPointer.Database(m.DatabaseName).Collection(collectionRecipes)
	if _, err := collection.InsertOne(ctx, &recipe); err != nil {
		return err
	}
	return nil
}
func (m *MongoRecipeRepository) GetAll(page, limit int) ([]Recipe, int64, error) {
	recipes := []Recipe{}
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64((page - 1) * limit))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := m.DbPointer.Database(m.DatabaseName).Collection(collectionRecipes)
	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	for cursor.Next(ctx) {
		var recipe Recipe
		if err = cursor.Decode(&recipe); err != nil {
			return nil, 0, err
		}
		recipes = append(recipes, recipe)
	}
	return recipes, 0, nil
}
func (m *MongoRecipeRepository) GetAllByTags(tags []int, page, limit int) ([]Recipe, int64, error) {
	recipes := []Recipe{}
	recipe := Recipe{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := m.DbPointer.Database(m.DatabaseName).Collection(collectionRecipes)
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64((page - 1) * limit))
	cursor, err := collection.Find(ctx, bson.M{"tags": bson.M{"$in": tags}}, findOptions)
	for cursor.Next(ctx) {
		if err = cursor.Decode(&recipe); err != nil {
			return nil, 0, err
		}
		recipes = append(recipes, recipe)
	}
	totalElements, err := collection.CountDocuments(ctx, bson.M{}, nil)
	if err != nil {
		return nil, 0, err
	}
	return recipes, totalElements, nil
}
