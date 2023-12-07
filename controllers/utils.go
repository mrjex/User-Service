package controllers

import (
	"Group20/Dentanoid/schemas"
	"encoding/json"
	"fmt"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
Extracts the topic from an mqtt message and returns the last token of its topic
e.g: '/users/mike' -> 'mike'
*/
type UpdateRequest struct {
	ID       primitive.ObjectID `json:"_id,omitempty"`
	OldName  string             `json:"oldname,omitempty"`
	Username string             `json:"username,omitempty"`
	Password string             `json:"password,omitempty"`
}

type Res struct {
	Status    int              `json:"status,omitempty"`
	RequestID string           `json:"requestID,omitempty"`
	Message   string           `json:"message,omitempty"`
	Patient   *schemas.Patient `json:"patient,omitempty"`
	Dentist   *schemas.Dentist `json:"dentist,omitempty"`
}

func GetPath(message mqtt.Message) string {
	tokens := strings.Split(message.Topic(), "/")
	result := tokens[len(tokens)-1]
	return result
}

// Adds mqtt code to stringified json
func AddCodeStringJson(json string, code string) string {
	var newJson string
	length := len(json)
	index := 0

	runes := []rune(json)

	for index >= 0 && index < (length-1) {
		newJson = newJson + string(runes[index])
		index++
	}
	newJson = newJson + ",\"Code\": \"" + code + "\"}"
	return newJson
}

func PublishReturnMessage(returnData Res, topic string, client mqtt.Client) {

	returnJson, err := json.Marshal(returnData)
	if err != nil {
		panic(err)
	}

	returnString := string(returnJson)
	fmt.Printf(returnString)

	client.Publish(topic, 0, false, returnString)

}
