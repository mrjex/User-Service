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
        go createPatient(payload.Username, payload.Password)
        fmt.Printf("%+v\n", payload)

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
    tokenUpdate := client.Subscribe("grp20/patient/update/+", byte(0), func(c mqtt.Client, m mqtt.Message) {

		var payload schemas.Patient
		username := GetPath(m)

		err := json.Unmarshal(m.Payload(), &payload)
		if err != nil {
			panic(err)
		}

		updatePatient(username, payload)
		fmt.Printf("%+v\n", payload)

	})

    if tokenUpdate.Error() != nil {
        panic(tokenRead.Error())
    }

    //REMOVE
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
    col := getPatientCollection()
    hashed, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    doc := schemas.Patient{Username: username, Password: string(hashed)}

    //if userExists(username) {
    //    return false;
    //}
    
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
    
    col := getPatientCollection()
    //Hash password
    hashed, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 14)

    // if userExists(payload.Username) {
        // return false
    //}

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


