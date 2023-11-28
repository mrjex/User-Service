package tests

import (
	"Group20/Dentanoid/controllers"
	"Group20/Dentanoid/database"
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
	code := m.Run()
	os.Exit(code)
}

func TestCreate(t *testing.T) {
	result := controllers.CreateDentist("mike", "password")
	if !result {
		t.Error("got no result")
	}

}
