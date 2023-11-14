package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func connectToDB() {

	c, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	client = c

	if err != nil {
		panic(err)
	}

	fmt.Println("App is connected to MongoDB")

}

func login(username string, password string) bool {

	data := User{}
	user := getUser(username)
	user.Decode(&data)

	fmt.Printf("Username: %s\n", data.Username)
	fmt.Printf("Hash: %s\n", data.Password)

	// Compare obtained password to the database hash
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(data.Password))
	if err == nil {
		fmt.Print("Authentication was successful")
		return true
	}
	fmt.Print("Authentication failed")
	return false

}

func signup(username string, password string) bool {

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

	fmt.Printf("Registered User ID: %v \n", result.InsertedID)
	return true

}

// func getUsers(w http.ResponseWriter, r *http.Request) {

// 	var users []User
// 	col := client.Database("dentanoid").Collection("users")
// 	doc, err := col.Find(context.TODO(), bson.D{{}})
// 	if err != nil {
// 		panic(err)
// 	}

// 	err = doc.All(context.TODO(), &users)
// 	if err != nil {
// 		panic(err)
// 	}

// 	result, err := json.Marshal(users)
// 	fmt.Printf("Sending users.\n")
// 	fmt.Fprint(w, string(result))

// }
