package controllers

import (
	"Group20/Dentanoid/database"
	"Group20/Dentanoid/schemas"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func InitialiseDentist(client mqtt.Client) {

	// 	CREATE
    tokenCreate := client.Subscribe("grp20/req/dentists/create", byte(0), func(c mqtt.Client, m mqtt.Message){

		var payload schemas.Dentist
		var returnData Res
        err1 := json.Unmarshal(m.Payload(), &payload)
        if err1 != nil {
            panic(err1)
        }
        err2 := json.Unmarshal(m.Payload(), &returnData)
        if err2 != nil {
            panic(err2)
        }
        go CreateDentist(payload, returnData, client)
    })
    if tokenCreate.Error() != nil {
        panic(tokenCreate.Error())
    }

	// READ
    tokenRead := client.Subscribe("grp20/req/dentists/read", byte(0), func(c mqtt.Client, m mqtt.Message){
        
        var payload schemas.Dentist
        var returnData Res

        err1 := json.Unmarshal(m.Payload(), &payload)
        if err1 != nil {
            panic(err1)
        }

        err2 := json.Unmarshal(m.Payload(), &returnData)
        if err2 != nil {
            panic(err2)
        }

        go GetDentist(payload.Username, returnData, client)
    })

    if tokenRead.Error() != nil {
        panic(tokenRead.Error())
    }

	// UPDATE
    tokenUpdate := client.Subscribe("grp20/req/dentists/update", byte(0), func(c mqtt.Client, m mqtt.Message) {


		var payload UpdateRequest
        var returnData Res


		err1 := json.Unmarshal(m.Payload(), &payload)
		if err1 != nil {
			panic(err1)
		}

        err2 := json.Unmarshal(m.Payload(), &returnData)
		if err2 != nil {
			panic(err2)
		}


		go UpdateDentist(payload, returnData, client)


	})

    if tokenUpdate.Error() != nil {
        panic(tokenRead.Error())
    }   

	//DELETE
    tokenRemove := client.Subscribe("grp20/req/dentists/delete", byte(0), func(c mqtt.Client, m mqtt.Message) {
        
        var payload schemas.Dentist
        var resData Res

        err1 := json.Unmarshal(m.Payload(), &payload)
        if err1 != nil {
            panic(err1)
        }

        err2 := json.Unmarshal(m.Payload(), &resData)
        if err2 != nil {
            panic(err2)
        }

        go DeleteDentist(payload.Username, resData, client)
    })

    if tokenRemove.Error() != nil{
        panic(tokenRemove.Error())
    }

}

// CREATE
func CreateDentist(dentist schemas.Dentist, returnData Res, client mqtt.Client) bool {

    var returnVal bool

    if userExists(dentist.Username) {
        returnData.Message = "User already exists"
        returnData.Status = 409
        returnVal = false
    }   else{

        col := getDentistCollection()
        hashed, err := bcrypt.GenerateFromPassword([]byte(dentist.Password), 12)
        doc := schemas.Dentist{Username: dentist.Username, Password: string(hashed)}

        dentist.Password = ""
        
        result, err := col.InsertOne(context.TODO(), doc)
        if err != nil {
            log.Fatal(err)
        }

        returnData.Message = "User created"
        returnData.Status = 201
        returnData.Dentist = &dentist

        fmt.Printf("Registered Dentist ID: %v \n", result.InsertedID)

        returnVal = true
    }

    PublishReturnMessage(returnData, "grp20/res/dentist/create", client)
    return returnVal

}

// READ
func GetDentist(id primitive.ObjectID, returnData Res, client mqtt.Client) bool {
    var returnVal bool

    col := getDentistCollection()
    user := &schemas.Dentist{}
    filter := bson.M{"_id": id}
    data := col.FindOne(context.TODO(), filter)
    data.Decode(user)

    if user.Username == ""{
        returnData.Message = "Dentist not found"
        returnData.Status = 404
        returnVal = false
    } else {
        returnVal = true
        returnData.Status = 200
        user.Password = ""
        returnData.Dentist = user
    }

    PublishReturnMessage(returnData, "grp20/res/dentists/read", client)

    return returnVal
}

// UPDATE
func UpdateDentist(payload UpdateRequest, returnData Res, client mqtt.Client) bool {
    var update bson.M
    var returnVal bool
    
    if userExists(payload.Username) {
        returnData.Message = "Username taken"
        returnData.Status = 409
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

        filter := bson.M{"_id": payload.ID}


        result, err := col.UpdateOne(context.TODO(), filter, update)

        if err != nil {
            log.Fatal(err)
            fmt.Printf("Updated failed for Dentist with Username: %v \n", payload.OldName)

            returnData.Status = 500
			returnData.Message = "Update failed"

            returnVal = false
        } else if result.MatchedCount == 1{
            fmt.Printf("Updated Dentist with Username: %v \n", payload.OldName)

            returnData.Status = 200
			returnData.Message = "Dentist updated"

            returnVal = true
        } else {
            fmt.Printf("No user with that name")

            returnData.Status = 404
			returnData.Message = "User not found"

            returnVal = false
        }

    }
        PublishReturnMessage(returnData, "grp20/res/dentist/update", client)
        return returnVal

}

// DELETE
func DeleteDentist(id primitive.ObjectID, returnData Res, client mqtt.Client) bool {
    var returnVal bool
    
    col := getDentistCollection()
    filter := bson.M{"_id": id}
    result, err := col.DeleteOne(context.TODO(), filter)


    if err != nil {
        log.Fatal(err)
    }

    if result.DeletedCount == 1 {

        returnData.Status = 200
		returnData.Message = "User with id: " + id.Hex() + " deleted"

        returnVal = true
	    fmt.Printf("Deleted Dentist %v \n", id.Hex())
    } else{

        returnData.Status = 404
		returnData.Message = "User not found"

        returnVal = false
    }

	PublishReturnMessage(returnData, "grp20/res/dentists/delete", client)
    return returnVal
}

func getDentistCollection() *mongo.Collection {
	col := database.Database.Collection("dentists")
	return col
}
