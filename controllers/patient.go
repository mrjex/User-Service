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
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func InitialisePatient (client mqtt.Client) {
    
    //CREATE
    tokenCreate := client.Subscribe("grp20/req/patient/create", byte(0), func(c mqtt.Client, m mqtt.Message){

		var payload schemas.Patient
        err := json.Unmarshal(m.Payload(), &payload)
        if err != nil {
            panic(err)
        }
        if createPatient(payload.Username, payload.Password) == true{
            fmt.Printf("%+v\n", payload)
        } else{
            fmt.Printf("Didnt work")
        }

    })
    if tokenCreate.Error() != nil {
        panic(tokenCreate.Error())
    }

    //READ
    tokenRead := client.Subscribe("grp20/req/patient/read", byte(0), func(c mqtt.Client, m mqtt.Message){
        
        var payload schemas.Patient
        err := json.Unmarshal(m.Payload(), &payload)
        if err != nil {
            panic(err)
        }
        user := getPatient(payload.Username)
		fmt.Printf("%+v\n", user)

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

		if updatePatient(username, payload) == true{
		    fmt.Printf("%+v\n", payload)
        } else{
            fmt.Printf("Didnt work")
        }


	})

    if tokenUpdate.Error() != nil {
        panic(tokenRead.Error())
    }   

    //REMOVE
    //TODO
    //Change subscription adress to get username in body

    tokenRemove := client.Subscribe("grp20/patient/delete/+", byte(0), func(c mqtt.Client, m mqtt.Message) {
        
        username := GetPath(m)
        deletePatient(username)
        fmt.Printf("Deleted Patient: %s", username)
    })

    if tokenRemove.Error() != nil{
        panic(tokenRemove.Error())
    }



}

//CREATE
func createPatient (username string, password string) bool {

    if userExists(username) {
        return false;
    }

    col := getPatientCollection()
    hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    doc := schemas.Patient{Username: username, Password: string(hashed)}
    
    result, err := col.InsertOne(context.TODO(), doc)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Registered Patient ID: %v \n", result.InsertedID)
    return true

}

//READ
func getPatient(username string) schemas.Patient {
    col := getPatientCollection()
    user := &schemas.Patient{}
    filter := bson.M{"username": username}
    data := col.FindOne(context.TODO(), filter)
    data.Decode(user)
    return *user
}

//UPDATE
func updatePatient(username string, payload schemas.Patient) bool {

    if userExists(payload.Username) {
        return false
    }
    
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

//REMOVE
func deletePatient(username string) bool {
    
    col := getPatientCollection()
    filter := bson.M{"username": username}
    result, err := col.DeleteOne(context.TODO(), filter)
    _ = result

    if err != nil {
        log.Fatal(err)
    }

	fmt.Printf("Deleted Patient: %v \n", username)
	return true
}

func getPatientCollection() *mongo.Collection {
    col := database.Database.Collection("Patients")
    return col
}


//TODO Responses
