package controllers

import (
	"Group20/Dentanoid/database"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username,omitempty"`
	Password string             `bson:"password,omitempty"`
}

func getUser(username string) *mongo.SingleResult {
	col := getUserCollection()
	user := col.FindOne(context.TODO(), User{Username: username})
	return user
}

func userExists(username string) bool {
	user := getUser(username)
	return user != nil
}

func getUserCollection() *mongo.Collection {
	col := database.Database.Collection("users")
	return col
}
