package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to MongoDB
	connectToDB()

	// Variant #1 MQTT
	/*
		In the likely scenario that REST will only be adhered to between the client and the main gateway; The communication between this module will
		most likely happen through a more efficient protocol like MQTT or HTTP 2.0. In that case, an MQTT client can be found below
		which can subscribe and publish to different topics and handle events
	*/
	mqttClient := getInstance()

	// CREATE
	mqttClient.Subscribe("/dentists/create", byte(0), func(c mqtt.Client, m mqtt.Message) {

		var payload User
		err := json.Unmarshal(m.Payload(), &payload)
		if err != nil {
			panic(err)
		}
		create(payload.Username, payload.Password)
		fmt.Printf("%+v\n", payload)

	})

	// READ
	mqttClient.Subscribe("/dentists/read/+", byte(0), func(c mqtt.Client, m mqtt.Message) {

		topic := strings.Split(m.Topic(), "/")
		username := topic[len(topic)-1]
		user := read(username)
		fmt.Printf("%+v\n", user)

	})

	// UPDATE
	mqttClient.Subscribe("/dentists/update/+", byte(0), func(c mqtt.Client, m mqtt.Message) {

		var payload User
		topic := strings.Split(m.Topic(), "/")
		username := topic[len(topic)-1]

		err := json.Unmarshal(m.Payload(), &payload)
		if err != nil {
			panic(err)
		}

		update(username, payload)
		fmt.Printf("%+v\n", payload)

	})

	//DELETE
	mqttClient.Subscribe("/dentists/delete/+", byte(0), func(c mqtt.Client, m mqtt.Message) {

		topic := strings.Split(m.Topic(), "/")
		username := topic[len(topic)-1]

		delete(username)
		fmt.Printf("Deleted Dentist: ", username)

	})

	// Variant #2 REST
	/*
		This endpoint listens for a request which contains a username and password field in its headers
		it then queries the database and proceeds to hash the given password to the one present in the database
		finally it responds with the status of the operation
	*/
	// http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
	// 	username := r.Header.Get("username")
	// 	password := r.Header.Get("password")
	// 	status := login(username, password)
	// 	fmt.Fprint(w, status)
	// })

	// http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {

	// 	username := r.Header.Get("username")
	// 	password := r.Header.Get("password")
	// 	status := signup(username, password)
	// 	fmt.Fprint(w, status)

	// })

	fmt.Println("HTTP server is listening for requests...")
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))

}
