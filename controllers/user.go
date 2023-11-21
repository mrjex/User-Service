package controllers

import (
	"Group20/Dentanoid/database"
	"Group20/Dentanoid/schemas"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username,omitempty"`
	Password string             `bson:"password,omitempty"`
}

func getUser(username string) *mongo.SingleResult {
	col := getUserCollection()
	user := col.FindOne(context.TODO(), User{Username: username})
	return user
}

func userExists(username string) bool {
    colPatients := getPatientCollection()
    filterPatients := bson.M{"username": username}
    patient := &schemas.Patient{}
    dataPatients := colPatients.FindOne(context.TODO(), filterPatients)
    dataPatients.Decode(patient)

    colDentists := getDentistCollection()
    filterDentists := bson.M{"username": username}
    dentist := &schemas.Dentist{}
    dataDentists := colDentists.FindOne(context.TODO(), filterDentists)
    dataDentists.Decode(dentist)


    if !(patient.Username == "" && dentist.Username == "") {
        fmt.Printf("There exists one")
    } else{
        fmt.Printf("There doesnt exist one")
    }


    return !(patient.Username == "" && dentist.Username == "")

}

func getUserCollection() *mongo.Collection {
	col := database.Database.Collection("users")
	return col
}
