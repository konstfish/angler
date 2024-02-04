package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id", json:"id"`
	Name     string             `bson:"name" json:"name"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"password"`
}

func (user *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) CheckPassword(providedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return false, err
	}
	return true, nil
}

/*
{
	"name": "konstfish",
	"email": "david@konst.fish"
	"password": "testpw1234"
}
*/
