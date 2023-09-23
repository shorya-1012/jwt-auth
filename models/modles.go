package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	UserId       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username     string             `json:"username" bson:"username"`
    Password     string             `json:"password" bson:"password"`
    RefreshToken string             `json:"-" bson:"refreshtoken"`
}

type ResponseUser struct {
	UserId       primitive.ObjectID `json:"_id,omitempty"` 
	Username     string             `json:"username"`
}
