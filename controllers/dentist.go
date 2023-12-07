package controllers

import (
	"Group20/Dentanoid/database"
	"Group20/Dentanoid/schemas"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"go.mongodb.org/mongo-driver/bson"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func InitialiseDentist(client mqtt.Client) {

	// 	CREATE
    tokenCreate := client.Subscribe("grp20/req/dentist/create", byte(0), func(c mqtt.Client, m mqtt.Message){

		var payload schemas.Dentist
        err := json.Unmarshal(m.Payload(), &payload)
        if err != nil {
            panic(err)
        }
        go CreateDentist(payload.Username, payload.Password, client)
    })
    if tokenCreate.Error() != nil {
        panic(tokenCreate.Error())
    }

	// READ
    tokenRead := client.Subscribe("grp20/req/dentist/read", byte(0), func(c mqtt.Client, m mqtt.Message){
        
        var payload schemas.Dentist
        err := json.Unmarshal(m.Payload(), &payload)
        if err != nil {
            panic(err)
        }
        go GetDentist(payload.Username, client)
    })

    if tokenRead.Error() != nil {
        panic(tokenRead.Error())
    }

	// UPDATE
    tokenUpdate := client.Subscribe("grp20/req/dentist/update", byte(0), func(c mqtt.Client, m mqtt.Message) {


		var payload UpdateRequest


		err := json.Unmarshal(m.Payload(), &payload)
		if err != nil {
			panic(err)
		}

		go UpdateDentist(payload, client)


	})

    if tokenUpdate.Error() != nil {
        panic(tokenRead.Error())
    }   

	//DELETE
    tokenRemove := client.Subscribe("grp20/req/dentist/delete", byte(0), func(c mqtt.Client, m mqtt.Message) {
        
        var payload schemas.Dentist
        err := json.Unmarshal(m.Payload(), &payload)
        if err != nil {
            panic(err)
        }

        go DeleteDentist(payload.Username, client)
    })

    if tokenRemove.Error() != nil{
        panic(tokenRemove.Error())
    }

}

// CREATE
func CreateDentist(username string, password string, client mqtt.Client) bool {

    var message string
    var returnVal bool

    if userExists(username) {
        message = "{\"Message\": \"User already exists\",\"Code\": \"409\"}"
        returnVal = false
    }   else{

        col := getDentistCollection()
        hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
        doc := schemas.Dentist{Username: username, Password: string(hashed)}
        
        result, err := col.InsertOne(context.TODO(), doc)
        if err != nil {
            log.Fatal(err)
        }

        message = "{\"Message\": \"User created\",\"Code\": \"201\"}"

        fmt.Printf("Registered Dentist ID: %v \n", result.InsertedID)

        returnVal = true
    }

    client.Publish("grp20/res/dentist/create", 0, false, message)
    return returnVal

}

// READ
func GetDentist(username string, client mqtt.Client) bool {
    var message string
    var code string
    var returnVal bool

    col := getDentistCollection()
    user := &schemas.Dentist{}
    filter := bson.M{"username": username}
    data := col.FindOne(context.TODO(), filter)
    data.Decode(user)

    jsonData, err := json.Marshal(user) 
    if err != nil{
        log.Fatal(err)
    }

    if user.Username == ""{
        code = "404"
        returnVal = false
    } else {
        code = "200"
        returnVal = true
    }
    message = AddCodeStringJson(string(jsonData), code)

    client.Publish("grp20/res/dentist/read", 0, false, message)

    return returnVal
}

// UPDATE
func UpdateDentist(payload UpdateRequest, client mqtt.Client) bool {
    var message string
    var code string
    var update bson.M
    var returnVal bool
    
    if userExists(payload.Username) {
        message = "{\"Message\": \"Username taken\"}"
        code = "409"
        returnVal = false
    } else{
        
        col := getDentistCollection()
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
            fmt.Printf("Updated failed for Dentist with Username: %v \n", payload.OldName)
            code = "500"
            message = "\"message\": \"Update failed\""
            returnVal = false
        } else if result.MatchedCount == 1{
            fmt.Printf("Updated Dentist with Username: %v \n", payload.OldName)
            code = "200"
            message = "\"message\": \"Dentist updated\""
            returnVal = true
        } else {
            fmt.Printf("No user with that name")
            code = "404"
            message = "\"message\": \"User not found\""
            returnVal = false
        }

    }
        message = AddCodeStringJson(message, code)
        client.Publish("grp20/res/dentist/update", 0, false, message)
        return returnVal

}

// DELETE
func DeleteDentist(username string, client mqtt.Client) bool {
    var message string
    var code string
    var returnVal bool
    
    col := getDentistCollection()
    filter := bson.M{"username": username}
    result, err := col.DeleteOne(context.TODO(), filter)


    if err != nil {
        log.Fatal(err)
    }

    if result.DeletedCount == 1 {
        message = "{\"Message\": \""+ username +" deleted\"}"
        code = "200"
        returnVal = true
	    fmt.Printf("Deleted Dentist %v \n", username)
    } else{
        message = "{\"Message\": \"Error deleting user\"}" 
        code = "404"
        returnVal = false
    }

    message = AddCodeStringJson(message, code)
    client.Publish("grp20/res/dentist/delete", 0, false, message)
    return returnVal
}

func getDentistCollection() *mongo.Collection {
	col := database.Database.Collection("dentists")
	return col
}
