package tests

import (
	"Group20/Dentanoid/controllers"
	"Group20/Dentanoid/mqtt"
	"Group20/Dentanoid/database"
	"testing"
    "github.com/joho/godotenv"
    "os"
    "log"

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
	controllers.DeletePatient("mike", client)
	controllers.DeletePatient("ben", client)
	result := controllers.CreatePatient("mike", "password", client)
	if !result {
		t.Error("Patient Creation Failed")
	}

    result = controllers.GetPatient("mike", client)
    if !result {
        t.Error("Reading patient failed")
    }

    updateUser := controllers.UpdateRequest{
        OldName: "mike",
        Username: "ben",
        Password: "123",
    }

    result = controllers.UpdatePatient(updateUser, client)
    if !result {
        t.Error("Updating patient failed")
    }

	result = controllers.DeletePatient("ben", client)
	if !result {
		t.Error("Patient Deletion Failed")
	}
}

