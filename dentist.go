package main

import (
	"context"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// CREATE
func create(username string, password string) bool {

	col := getUserCollection()

	// Hash the password using Bcrypt
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	doc := User{Username: username, Password: string(hashed)}

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
func read(username string) User {
	col := getUserCollection()
	data := col.FindOne(context.TODO(), User{Username: username})
	user := User{}
	data.Decode(user)
	return user
}

// UPDATE
func update(username string, payload User) bool {

	col := getUserCollection()

	// Hash the password using Bcrypt
	hashed, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 14)
	doc := User{Username: payload.Username, Password: string(hashed)}

	if userExists(payload.Username) {
		return false
	}

	result, err := col.UpdateOne(context.TODO(), User{Username: payload.Username}, doc)
	_ = result

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Updated Dentist with Username: %v \n", username)
	return true

}

// DELETE
func delete(username string) bool {

	col := getUserCollection()

	result, err := col.DeleteOne(context.TODO(), User{Username: username})
	_ = result

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Deleted Dentist: %v \n", username)
	return true

}
