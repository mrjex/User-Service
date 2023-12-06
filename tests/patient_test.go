package tests

import (
	"Group20/Dentanoid/controllers"
	"Group20/Dentanoid/database"
	"Group20/Dentanoid/mqtt"
	"Group20/Dentanoid/schemas"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var client MQTT.Client

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Error loading .env file")
	}
	database.Connect()
	client = mqtt.GetInstance()
	code := m.Run()
	os.Exit(code)
}


func TestPatientCRUD(t *testing.T) {

    var res controllers.Res

    payload := schemas.Patient{
        ID: primitive.NewObjectID(),
        Username: "mike",
        Password: "password",
    }

    id := payload.ID

	result := controllers.CreatePatient(payload, res, client)
	if !result {
		t.Error("Patient Creation Failed")
	}

    result = controllers.GetPatient(id, res, client)
    if !result {
        t.Error("Reading patient failed")
    }

    updateUser := controllers.UpdateRequest{
        ID: payload.ID,
        OldName: "mike",
        Username: "ben",
        Password: "123",
    }

    result = controllers.UpdatePatient(updateUser, res, client)
    if !result {
        t.Error("Updating patient failed")
    }

	result = controllers.DeletePatient(updateUser.ID, res, client)
	if !result {
		t.Error("Patient Deletion Failed")
	}
}

