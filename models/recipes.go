package models

import (
	"context"
	"fmt"
	"time"

	"github.com/Daniorocket/cookit-backend/lib"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collectionRecipes = "recipes"
var collectionUnits = "units"

type Unit struct {
	ID     string `json:"id" bson:"id"`
	Name   string `json:"name" bson:"name"`
	Symbol string `json:"symbol" bson:"symbol"`
}

type Ingredient struct {
	ID     string `json:"id" bson:"id" validate:"nonnil,nonzero"`
	Name   string `json:"name" bson:"name" validate:"nonnil,nonzero"`
	Count  int    `json:"count" bson:"count" validate:"nonnil,nonzero"`
	UnitID string `json:"unitID" bson:"unit_id" validate:"nonnil,nonzero"`
}

type Recipe struct {
	ID            string       `json:"id" bson:"id"`
	Name          string       `json:"name" bson:"name" validate:"nonnil,nonzero"`
	Username      string       `json:"username" bson:"username" validate:"nonnil,nonzero"`
	Difficulty    int          `json:"difficulty" bson:"difficulty" validate:"min=1,max=3"`
	Ingredients   []Ingredient `json:"ingredients" bson:"ingredients" validate:"nonnil,nonzero,min=1"`
	Steps         []string     `json:"steps" bson:"steps" validate:"nonnil,nonzero,min=1"`
	CategoriesID  []string     `json:"categoriesID" bson:"categories_id" validate:"nonnil,nonzero,min=1"`
	Description   string       `json:"description" bson:"description" validate:"nonnil,nonzero"`
	Date          string       `json:"date" bson:"date" validate:"nonnil,nonzero"`
	PreparingTime int          `json:"preparingTime" bson:"preparing_time" validate:"nonnil,nonzero"`
	File          lib.File     `json:"file" bson:"file" validate:"nonnil,nonzero"`
}

type MongoRecipeRepository struct {
	DbPointer    *mongo.Client
	DatabaseName string
}

type PostgreSQLRecipeRepository struct {
	ConStr       string
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
func (m *MongoRecipeRepository) GetAll(categories []string, page, limit int, name string) ([]Recipe, int64, error) {
	var filter primitive.M
	recipes := []Recipe{}
	recipe := Recipe{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := m.DbPointer.Database(m.DatabaseName).Collection(collectionRecipes)
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64((page - 1) * limit))
	if categories[0] != "" && len(name) > 0 {
		filter = bson.M{
			"$and": []bson.M{
				{"categories_id": bson.M{"$in": categories}},
				{"name": primitive.Regex{Pattern: name, Options: "i"}},
			},
		}
	} else if categories[0] != "" {
		filter = bson.M{"categories_id": bson.M{"$in": categories}}
	} else if len(name) > 0 {
		filter = bson.M{"name": primitive.Regex{Pattern: name, Options: ""}}
	}
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	for cursor.Next(ctx) {
		if err = cursor.Decode(&recipe); err != nil {
			return nil, 0, err
		}
		recipes = append(recipes, recipe)
	}
	totalElements, err := collection.CountDocuments(ctx, filter, nil)
	if err != nil {
		return nil, 0, err
	}
	return recipes, totalElements, nil
}
func (m *MongoRecipeRepository) GetUnit(unitID string) (Unit, error) {
	unit := Unit{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := m.DbPointer.Database(m.DatabaseName).Collection(collectionUnits)

	if err := collection.FindOne(ctx, bson.M{"id": unitID}).Decode(&unit); err != nil {
		return Unit{}, err
	}
	return unit, nil
}
func (m *MongoRecipeRepository) GetAllUnits() ([]Unit, error) {
	unit := Unit{}
	units := []Unit{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := m.DbPointer.Database(m.DatabaseName).Collection(collectionUnits)
	cursor, err := collection.Find(ctx, bson.M{}, nil)
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		if err = cursor.Decode(&unit); err != nil {
			return nil, err
		}
		units = append(units, unit)
	}
	return units, nil
}
func (d *MongoRecipeRepository) GetByID(id string) (Recipe, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := d.DbPointer.Database(d.DatabaseName).Collection("recipes")
	rec := Recipe{}
	if err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&rec); err != nil {
		return Recipe{}, err
	}
	return rec, nil
}
func (d *MongoRecipeRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := d.DbPointer.Database(d.DatabaseName).Collection("recipes")
	result, err := collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}
	fmt.Printf("deleted  %v document(s)\n", result.DeletedCount)
	return nil
}
func (d *MongoRecipeRepository) GetByUser(username string, page, limit int) ([]Recipe, int64, error) {
	var recipe Recipe
	recipes := []Recipe{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := d.DbPointer.Database(d.DatabaseName).Collection(collectionRecipes)
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64((page - 1) * limit))
	cursor, err := collection.Find(ctx, bson.M{"username": username}, findOptions)
	if err != nil {
		return nil, 0, err
	}
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

func (p *PostgreSQLRecipeRepository) Create(recipe Recipe) error {
	fmt.Println("Not implemented yet")
	return nil
}

func (p *PostgreSQLRecipeRepository) Delete(id string) error {
	fmt.Println("Not implemented yet")
	return nil
}

func (p *PostgreSQLRecipeRepository) GetAll(categories []string, page, limit int, name string) ([]Recipe, int64, error) {
	fmt.Println("Not implemented yet")
	return nil, 0, nil
}

func (p *PostgreSQLRecipeRepository) GetAllUnits() ([]Unit, error) {
	fmt.Println("Not implemented yet")
	return nil, nil
}

func (p *PostgreSQLRecipeRepository) GetByID(id string) (Recipe, error) {
	fmt.Println("Not implemented yet")
	return Recipe{}, nil
}

func (p *PostgreSQLRecipeRepository) GetByUser(username string, page, limit int) ([]Recipe, int64, error) {
	fmt.Println("Not implemented yet")
	return nil, 0, nil
}

func (p *PostgreSQLRecipeRepository) GetUnit(unitID string) (Unit, error) {
	fmt.Println("Not implemented yet")
	return Unit{}, nil
}
