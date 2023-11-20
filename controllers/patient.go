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
	"golang.org/x/crypto/bcrypt"
)

func InitialisePatient (client mqtt.Client) {
    
    token := client.Subscribe("grp20/req/patient/create", byte(0), func(c mqtt.Client, m mqtt.Message){

		var payload schemas.Patient
        err := json.Unmarshal(m.Payload(), &payload)
        if err != nil {
            panic(err)
        }
        go createPatient(payload.Username, payload.Password)
        fmt.Printf("%+v\n", payload)

    })
    if token.Error() != nil {
        panic(token.Error())
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

func getPatientCollection() *mongo.Collection {
    col := database.Database.Collection("Patients")
    return col
}


