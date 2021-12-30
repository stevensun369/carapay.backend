package models

import (
	"backend/db"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
  ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID       string   `json:"userID,omitempty" bson:"userID,omitempty"`
	Username     string   `json:"username,omitempty" bson:"username,omitempty"`
	Email string `json:"email,omitempty" bson:"email,omitempty"`
  Password     string   `json:"password,omitempty" bson:"password,omitempty"`
	People []string `json:"people,omitempty" bson:"people,omitempty"`
  CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
  UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

func GetUserByID(userID string) (User) {
	var user User
	usersCollection := db.GetCollection("users")

	usersCollection.FindOne(context.Background(), bson.M{
		"userID": userID,
	}).Decode(&user)

	return user
}