package tests

import (
	"Group20/Dentanoid/controllers"
	"Group20/Dentanoid/schemas"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDentistCRUD(t *testing.T) {

	var res controllers.Res

	payload := schemas.Dentist{
		ID:       primitive.NewObjectID(),
		Username: "mike",
		Password: "password",
	}

	id := payload.ID

	result := controllers.CreateDentist(payload, res, client)
	if !result {
		t.Error("Dentist Creation Failed")
	}

	result = controllers.GetDentist(id, res, client)
	if !result {
		t.Error("Reading patient failed")
	}

	updateUser := controllers.UpdateRequest{
		ID:       payload.ID,
		OldName:  "mike",
		Username: "ben",
		Password: "123",
	}

	result = controllers.UpdateDentist(updateUser, res, client)
	if !result {
		t.Error("Updating patient failed")
	}

	result = controllers.DeleteDentist(updateUser.ID, res, client)
	if !result {
		t.Error("Dentist Deletion Failed")
	}
}
