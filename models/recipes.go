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

type Recipe struct {
	ID          string
	Name        string
	UserID      string
	Kitchen     KitchenStyle
	Tags        Tags
	ListOfSteps []string
	Description string
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
