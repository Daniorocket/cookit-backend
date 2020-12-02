package models

import (
	"errors"

	"github.com/globalsign/mgo"
)

var CollectionUsers = "users"
var ErrCreateUser = errors.New("Failed to create user record")

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

func RegisterUser(session *mgo.Session, db string, u *User) error {

	sess := session.Copy()
	defer sess.Close()
	if err := sess.DB(db).C(CollectionUsers).Insert(&u); err != nil {
		return ErrCreateUser
	}
	return nil
}
