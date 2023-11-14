package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username,omitempty"`
	Password string             `bson:"password,omitempty"`
}

var client *mongo.Client

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connectToDB()
	http.HandleFunc("/api/", handler)
	http.HandleFunc("/users", getUsers)
	fmt.Println("Server is listening for connections...")
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}

func getUsers(w http.ResponseWriter, r *http.Request) {

	var users []User
	col := client.Database("dentanoid").Collection("users")
	doc, err := col.Find(context.TODO(), bson.D{{}})
	if err != nil {
		panic(err)
	}

	err = doc.All(context.TODO(), &users)
	if err != nil {
		panic(err)
	}

	result, err := json.Marshal(users)
	fmt.Printf("Sending users.\n")
	fmt.Fprint(w, string(result))

}

func handler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Path[1:]
	addUser(username, "password")
}

func connectToDB() {

	c, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	client = c

	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to MongoDB")

}

func addUser(username string, password string) {

	col := client.Database("dentanoid").Collection("users")

	if col == nil {
		fmt.Println("Collection does not exist")
	}

	doc := User{Username: username, Password: password}

	result, err := col.InsertOne(context.TODO(), doc)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Result ID: %v \n", result.InsertedID)
	fmt.Printf("Added User with username %s and password %s\n", username, password)

}
