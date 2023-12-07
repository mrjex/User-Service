package tests

import (
	"Group20/Dentanoid/controllers"
	"Group20/Dentanoid/schemas"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreate(t *testing.T) {
    dentist := schemas.Dentist{
        Username: "mike",
        Password: "password",
        ID: primitive.NewObjectID(),
    }

    var returnData controllers.Res

	result := controllers.CreateDentist(dentist, returnData, client)
	if !result {
		t.Error("Dentist Creation Failed")
	}
	result = controllers.DeleteDentist(dentist.Username, returnData, client)
	if !result {
		t.Error("Dentist Deletion Failed")
	}
}
