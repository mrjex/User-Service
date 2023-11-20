package controllers

import (
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

/*
Extracts the topic from an mqtt message and returns the last token of its topic
e.g: '/users/mike' -> 'mike'
*/
func GetPath(message mqtt.Message) string {
	tokens := strings.Split(message.Topic(), "/")
	result := tokens[len(tokens)-1]
	return result
}
