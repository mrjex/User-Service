package controllers

import (
	"Group20/Dentanoid/database"
	"Group20/Dentanoid/schemas"
	"context"
	"encoding/json"

	//"encoding/json"
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func InitialisePatient(client mqtt.Client) {

	//CREATE
	tokenCreate := client.Subscribe("grp20/req/patient/create", byte(0), func(c mqtt.Client, m mqtt.Message) {

		var payload schemas.Patient
		err := json.Unmarshal(m.Payload(), &payload)
		if err != nil {
			panic(err)
		}
		go CreatePatient(payload.Username, payload.Password, client)
	})
	if tokenCreate.Error() != nil {
		panic(tokenCreate.Error())
	}

	//READ
	tokenRead := client.Subscribe("grp20/req/patient/read", byte(0), func(c mqtt.Client, m mqtt.Message) {

		var payload schemas.Patient
		err := json.Unmarshal(m.Payload(), &payload)
		if err != nil {
			panic(err)
		}
		go getPatient(payload.Username, client)
	})

	if tokenRead.Error() != nil {
		panic(tokenRead.Error())
	}

	//UPDATE
	//TODO
	//Change subscription adress to get username in body
	tokenUpdate := client.Subscribe("grp20/patient/update/+", byte(0), func(c mqtt.Client, m mqtt.Message) {

		var payload schemas.Patient
		username := GetPath(m)

		err := json.Unmarshal(m.Payload(), &payload)
		if err != nil {
			panic(err)
		}

		if updatePatient(username, payload) {
			fmt.Printf("%+v\n", payload)
		} else {
			fmt.Printf("Didnt work")
		}

	})

	if tokenUpdate.Error() != nil {
		panic(tokenRead.Error())
	}

	//REMOVE
	//TODO
	//Change subscription adress to get username in body

	tokenRemove := client.Subscribe("grp20/req/patient/delete", byte(0), func(c mqtt.Client, m mqtt.Message) {

		var payload schemas.Patient
		err := json.Unmarshal(m.Payload(), &payload)
		if err != nil {
			panic(err)
		}

		go DeletePatient(payload.Username, client)
	})

	if tokenRemove.Error() != nil {
		panic(tokenRemove.Error())
	}

}

// CREATE
func CreatePatient(username string, password string, client mqtt.Client) bool {
	var message string

	if userExists(username) {
		message = "{\"Message\": \"User already exists\",\"Code\": \"409\"}"
		return false
	}

	col := getPatientCollection()
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	doc := schemas.Patient{Username: username, Password: string(hashed)}

	result, err := col.InsertOne(context.TODO(), doc)
	if err != nil {
		log.Fatal(err)
	}
	message = "{\"Message\": \"User created\",\"Code\": \"201\"}"
	fmt.Printf("Registered Patient ID: %v \n", result.InsertedID)
	client.Publish("grp20/res/patient/create", 0, false, message)
	return true

}

// READ
func getPatient(username string, client mqtt.Client) {
	var message string
	var code string

	col := getPatientCollection()
	user := &schemas.Patient{}
	filter := bson.M{"username": username}
	data := col.FindOne(context.TODO(), filter)
	data.Decode(user)

	jsonData, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	if user.Username == "" {
		code = "404"
	} else {
		code = "200"
		fmt.Printf(user.Username)
	}
	message = AddCodeStringJson(string(jsonData), code)

	client.Publish("grp20/res/patient/read", 0, false, message)

}

// UPDATE
func updatePatient(username string, payload schemas.Patient) bool {

	// if userExists(payload.Username) {
	//     message := "{\"Message\": \"Username taken\"}"
	//     code := "409"
	// }

	col := getPatientCollection()
	//Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 12)

	update := bson.M{"$set": bson.M{"username": payload.Username, "password": string(hashed)}}
	filter := bson.M{"username": username}

	result, err := col.UpdateOne(context.TODO(), filter, update)
	_ = result

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Updated Patient with Username: %v \n", username)
	return true

}

// REMOVE
func DeletePatient(username string, client mqtt.Client) bool {
	var message string
	var code string

	col := getPatientCollection()
	filter := bson.M{"username": username}
	result, err := col.DeleteOne(context.TODO(), filter)

	if err != nil {
		log.Fatal(err)
	}

	if result.DeletedCount == 0 {
		message = "{\"Message\": \"Error deleting user\"}"
		code = "404"
	} else {
		fmt.Printf("Deleted Patient: %v \n", username)
		message = "{\"Message\": \"" + username + " deleted\"}"
		code = "200"
	}

	message = AddCodeStringJson(message, code)
	client.Publish("grp20/res/patient/delete", 0, false, message)
	return code == "200"

}

func getPatientCollection() *mongo.Collection {
	col := database.Database.Collection("Patients")
	return col
}

//TODO Responses
