package controllers

import (
	"Group20/Dentanoid/database"
	"Group20/Dentanoid/schemas"
	"context"
	"encoding/json"
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func InitialiseDentist(client mqtt.Client) {

	// 	CREATE
	client.Subscribe("grp20/dentists/post", byte(0), func(c mqtt.Client, m mqtt.Message) {

		var payload User
		err := json.Unmarshal(m.Payload(), &payload)
		if err != nil {
			panic(err)
		}
		CreateDentist(payload.Username, payload.Password)
		fmt.Printf("%+v\n", payload)

	})

	// READ
	client.Subscribe("grp20/dentists/get/+", byte(0), func(c mqtt.Client, m mqtt.Message) {

		username := GetPath(m)
		user := getDentist(username)
		fmt.Printf("%+v\n", user)

	})

	// UPDATE
	client.Subscribe("grp20/dentists/update/+", byte(0), func(c mqtt.Client, m mqtt.Message) {

		var payload schemas.Dentist
		username := GetPath(m)

		err := json.Unmarshal(m.Payload(), &payload)
		if err != nil {
			panic(err)
		}

		updateDentist(username, payload)
		fmt.Printf("%+v\n", payload)

	})

	//DELETE
	client.Subscribe("grp20/dentists/delete/+", byte(0), func(c mqtt.Client, m mqtt.Message) {

		username := GetPath(m)
		deleteDentist(username)
		fmt.Printf("Deleted Dentist: %s", username)

	})

}

// CREATE
func CreateDentist(username string, password string) bool {

	col := getDentistCollection()
	// Hash the password using Bcrypt
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	doc := schemas.Dentist{Username: username, Password: string(hashed)}

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
func getDentist(username string) schemas.Dentist {
	col := getDentistCollection()
	data := col.FindOne(context.TODO(), schemas.Dentist{Username: username})
	user := schemas.Dentist{}
	data.Decode(user)
	return user
}

// UPDATE
func updateDentist(username string, payload schemas.Dentist) bool {

	col := getDentistCollection()
	// Hash the password using Bcrypt
	hashed, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 14)
	doc := schemas.Dentist{Username: payload.Username, Password: string(hashed)}

	if userExists(payload.Username) {
		return false
	}

	result, err := col.UpdateOne(context.TODO(), schemas.Dentist{Username: payload.Username}, doc)
	_ = result

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Updated Dentist with Username: %v \n", username)
	return true

}

// DELETE
func deleteDentist(username string) bool {

	col := getDentistCollection()
	result, err := col.DeleteOne(context.TODO(), schemas.Dentist{Username: username})
	_ = result

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Deleted Dentist: %v \n", username)
	return true

}

func getDentistCollection() *mongo.Collection {
	col := database.Database.Collection("dentists")
	return col
}
