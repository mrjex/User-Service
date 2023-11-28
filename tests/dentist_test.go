package tests

import (
	"Group20/Dentanoid/controllers"
	"Group20/Dentanoid/database"
	"Group20/Dentanoid/mqtt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	database.Connect()
	client = mqtt.GetInstance()
	code := m.Run()
	os.Exit(code)
}

func TestCreate(t *testing.T) {
	result := controllers.CreateDentist("mike", "password")
	if !result {
		t.Error("Dentist Creation Failed")
	}
	result = controllers.DeleteDentist("mike")
	if !result {
		t.Error("Dentist Deletion Failed")
	}
}
