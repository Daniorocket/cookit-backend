package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Daniorocket/cookit-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CategoryRepository interface {
	GetByID(id string) (models.Category, error)
	Create(category models.Category) error
	GetAll(page, limit int) ([]models.Category, int64, error)
}

type RecipeRepository interface {
	Create(recipe models.Recipe) error
	GetAll(page, limit int) ([]models.Recipe, int64, error)
	GetAllByTags(tags []int, page, limit int) ([]models.Recipe, int64, error)
}

type AuthRepository interface {
	Register(user models.User) error
	GetPassword(login string) (string, error)
	GetUserinfo(username string) (models.User, error)
}

func InitMongoDatabase() (CategoryRepository, RecipeRepository, AuthRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		os.Getenv("MONGODB_URI"),
	))
	if err != nil {
		log.Println("Err connect mongo:", err)
		return nil, nil, nil, err
	}

	categoryRepository := &models.MongoCategoryRepository{
		DbPointer:    client,
		DatabaseName: os.Getenv("DBName"),
	}
	recipeRepository := &models.MongoRecipeRepository{
		DbPointer:    client,
		DatabaseName: os.Getenv("DBName"),
	}
	authRepository := &models.MongoAuthRepository{
		DbPointer:    client,
		DatabaseName: os.Getenv("DBName"),
	}
	collection := client.Database("CookIt").Collection("users")
	//Index for users
	keys := []string{"id", "email", "username"}
	for i := range keys {
		if _, err := collection.Indexes().CreateOne(
			context.Background(),
			mongo.IndexModel{
				Keys: bson.M{
					keys[i]: 1,
				},
				Options: options.Index().SetUnique(true).SetBackground(true).SetSparse(true),
			},
		); err != nil {
			return nil, nil, nil, err
		}
	}
	collection = client.Database("CookIt").Collection("recipes")
	//Index for recipes
	keys = []string{"id"}
	for i := range keys {
		if _, err := collection.Indexes().CreateOne(
			context.Background(),
			mongo.IndexModel{
				Keys: bson.M{
					keys[i]: 1,
				},
				Options: options.Index().SetUnique(true).SetBackground(true).SetSparse(true),
			},
		); err != nil {
			return nil, nil, nil, err
		}
	}
	collection = client.Database("CookIt").Collection("categories")
	//Index for categories
	keys = []string{"id", "label_pl", "label_en"}
	for i := range keys {
		if _, err := collection.Indexes().CreateOne(
			context.Background(),
			mongo.IndexModel{
				Keys: bson.M{
					keys[i]: 1,
				},
				Options: options.Index().SetUnique(true).SetBackground(true).SetSparse(true),
			},
		); err != nil {
			return nil, nil, nil, err
		}
	}
	return categoryRepository, recipeRepository, authRepository, nil
}
