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

func InitialisePatient (client mqtt.Client) {
    
    //CREATE
    tokenCreate := client.Subscribe("grp20/req/patient/create", byte(0), func(c mqtt.Client, m mqtt.Message){

		var payload schemas.Patient
        err := json.Unmarshal(m.Payload(), &payload)
        if err != nil {
            panic(err)
        }
        go createPatient(payload.Username, payload.Password, client)
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
        go getPatient(payload.Username, client)
    })

    if tokenRead.Error() != nil {
        panic(tokenRead.Error())
    }


    //UPDATE
    //TODO
    //Change subscription adress to get username in body
    tokenUpdate := client.Subscribe("grp20/req/patient/update", byte(0), func(c mqtt.Client, m mqtt.Message) {


		var payload updateRequest


		err := json.Unmarshal(m.Payload(), &payload)
		if err != nil {
			panic(err)
		}

		go updatePatient(payload, client)


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

        go deletePatient(payload.Username, client)
    })

    if tokenRemove.Error() != nil{
        panic(tokenRemove.Error())
    }



}

//CREATE
func createPatient (username string, password string, client mqtt.Client) {
    var message string

    if userExists(username) {
        message = "{\"Message\": \"User already exists\",\"Code\": \"409\"}"
    }   else{

        col := getPatientCollection()
        hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
        doc := schemas.Patient{Username: username, Password: string(hashed)}
        
        result, err := col.InsertOne(context.TODO(), doc)
        if err != nil {
            log.Fatal(err)
        }

        message = "{\"Message\": \"User created\",\"Code\": \"201\"}"

        fmt.Printf("Registered Patient ID: %v \n", result.InsertedID)
    }

    client.Publish("grp20/res/patient/create", 0, false, message)
}

//READ
func getPatient(username string, client mqtt.Client){
    var message string
    var code string

    col := getPatientCollection()
    user := &schemas.Patient{}
    filter := bson.M{"username": username}
    data := col.FindOne(context.TODO(), filter)
    data.Decode(user)

    jsonData, err := json.Marshal(user) 
    if err != nil{
        log.Fatal(err)
    }

    if user.Username == ""{
        code = "404"
    } else {
        code = "200"
        fmt.Printf(user.Username)
    }
    message = AddCodeStringJson(string(jsonData), code)

    client.Publish("grp20/res/patient/read", 0, false, message)


}

//UPDATE
func updatePatient(payload updateRequest, client mqtt.Client) {
    var message string
    var code string
    var update bson.M
    
    fmt.Printf(string(payload.Username))


    if userExists(payload.Username) {
        message = "{\"Message\": \"Username taken\"}"
        code = "409"
    } else{
        
        col := getPatientCollection()
        //Hash password, might introduce performance issues when done before checking if olduser exists
        hashed, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 12)

        if (payload.Username != "") && (payload.Password != ""){
            update = bson.M{"$set": bson.M{"username": payload.Username, "password": string(hashed)}}
        } else if payload.Username != "" {
            update = bson.M{"$set": bson.M{"username": payload.Username}}
        } else if payload.Password != ""{
            update = bson.M{"$set": bson.M{"password": string(hashed)}}
        }

        filter := bson.M{"username": payload.OldName}


        result, err := col.UpdateOne(context.TODO(), filter, update)

        if err != nil {
            log.Fatal(err)
            fmt.Printf("Updated failed for Patient with Username: %v \n", payload.OldName)
            code = "500"
            message = "\"message\": \"Update failed\""
        } else if result.MatchedCount == 1{
            fmt.Printf("Updated Patient with Username: %v \n", payload.OldName)
            code = "200"
            message = "\"message\": \"Patient updated\""
        } else {
            fmt.Printf("No user with that name")
            code = "404"
            message = "\"message\": \"User not found\""
        }

    }
        message = AddCodeStringJson(message, code)
        client.Publish("grp20/res/patient/update", 0, false, message)
}

//REMOVE
func deletePatient(username string, client mqtt.Client) {
    var message string
    var code string
    
    col := getPatientCollection()
    filter := bson.M{"username": username}
    result, err := col.DeleteOne(context.TODO(), filter)


    if err != nil {
        log.Fatal(err)
    }

    if result.DeletedCount == 1 {
        message = "{\"Message\": \""+ username +" deleted\"}"
        code = "200"
	    fmt.Printf("Deleted Patient: %v \n", username)
    } else{
        message = "{\"Message\": \"Error deleting user\"}" 
        code = "404"
    }

    message = AddCodeStringJson(message, code)
    client.Publish("grp20/res/patient/delete", 0, false, message)
}

func getPatientCollection() *mongo.Collection {
    col := database.Database.Collection("Patients")
    return col
}


//TODO Responses
