package models

import (
	"errors"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

var CollectionUsers = "users"
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

func RegisterUser(session *mgo.Session, db string, u *User) error {

	sess := session.Copy()
	defer sess.Close()
	if err := sess.DB(db).C(CollectionUsers).Insert(&u); err != nil {
		return ErrCreateUser
	}
	return nil
}
func GetPasswordByUsernameOrEmail(session *mgo.Session, db string, username string) (string, error) {
	sess := session.Copy()
	defer sess.Close()
	result := User{}
	if err := sess.DB(db).C(CollectionUsers).Find(bson.M{"username": username}).One(&result); err == nil {
		return result.Password, nil
	}
	if err := sess.DB(db).C(CollectionUsers).Find(bson.M{"email": username}).One(&result); err == nil {
		return result.Password, nil
	}
	return "", ErrFindUser
}
