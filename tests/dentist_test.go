package tests

import (
	"Group20/Dentanoid/controllers"
	"testing"
)

func TestCreate(t *testing.T) {
	result := controllers.CreateDentist("mike", "password", client)
	if !result {
		t.Error("Dentist Creation Failed")
	}
	result = controllers.DeleteDentist("mike", client)
	if !result {
		t.Error("Dentist Deletion Failed")
	}
}
