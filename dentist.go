package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type Dentist struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username,omitempty"`
	Password string             `bson:"password,omitempty"`
}

// CREATE
func create(username string, password string) bool {

	col := getDentistCollection()
	// Hash the password using Bcrypt
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	doc := Dentist{Username: username, Password: string(hashed)}

	if userExists(username) {
		return false
	}

	result, err := col.InsertOne(context.TODO(), doc)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Registered Dentist ID: %v \n", result.InsertedID)
	return true

}

// READ
func read(username string) Dentist {
	col := getDentistCollection()
	data := col.FindOne(context.TODO(), Dentist{Username: username})
	user := Dentist{}
	data.Decode(user)
	return user
}

// UPDATE
func update(username string, payload Dentist) bool {

	col := getDentistCollection()
	// Hash the password using Bcrypt
	hashed, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 14)
	doc := Dentist{Username: payload.Username, Password: string(hashed)}

	if userExists(payload.Username) {
		return false
	}

	result, err := col.UpdateOne(context.TODO(), Dentist{Username: payload.Username}, doc)
	_ = result

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Updated Dentist with Username: %v \n", username)
	return true

}

// DELETE
func delete(username string) bool {

	col := getDentistCollection()
	result, err := col.DeleteOne(context.TODO(), Dentist{Username: username})
	_ = result

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Deleted Dentist: %v \n", username)
	return true

}
