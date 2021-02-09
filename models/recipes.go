package models

import (
	"errors"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

var CollectionRecipes = "recipes"
var ErrCreateRecipe = errors.New("Failed to create recipe record")
var ErrFindRecipe = errors.New("Failed to find recipe record")

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

type Categories int

const (
	Category1 Categories = iota
	Category2
	Category3
)

type Recipe struct {
	ID               string       `json:"id" bson:"id"`
	Name             string       `json:"name" bson:"name"`
	UserID           string       `json:"userId" bson:"user_id"`
	Kitchen          KitchenStyle `json:"kitchenStyle" bson:"kitchen_style"`
	Tags             Tags         `json:"tags" bson:"tags"`
	ListOfSteps      []string     `json:"listOfSteps" bson:"list_of_steps"`
	ListOfCategories []Categories `json:"categories" bson:"categories"`
	Description      string       `json:"description" bson:"description"`
}

func CreateRecipe(session *mgo.Session, db string, r Recipe) error {
	sess := session.Copy()
	defer sess.Close()
	if err := sess.DB(db).C(CollectionRecipes).Insert(&r); err != nil {
		return ErrCreateRecipe
	}
	return nil
}
func GetAllRecipes(session *mgo.Session, db string) ([]Recipe, error) {
	recipes := []Recipe{}
	sess := session.Copy()
	defer sess.Close()
	if err := sess.DB(db).C(CollectionRecipes).Find(bson.M{}).All(&recipes); err != nil {
		return nil, ErrFindRecipe
	}
	return recipes, nil
}
