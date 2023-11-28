package tests

import (
	"Group20/Dentanoid/controllers"
	"Group20/Dentanoid/mqtt"
	"testing"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var client MQTT.Client

func TestPatientCreate(t *testing.T) {
	client := mqtt.GetInstance()
	controllers.DeletePatient("mike", client)
	result := controllers.CreatePatient("mike", "password", client)
	if !result {
		t.Error("Patient Creation Failed")
	}
	result = controllers.DeletePatient("mike", client)
	if !result {
		t.Error("Dentist Deletion Failed")
	}
}
