package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Battlertrunks/api"
	"github.com/Battlertrunks/wsMessenger"
)

func setupRoutes() {
	// Initialize the REST API
	err := api.RunGin()
	if err != nil {
		log.Fatal("Database Failed")
		return
	}

	// Initialize the WebSocket connection
	// TODO: We should start the websocket based on when the user is in a live chat connection with
	// a another receiver
	http.HandleFunc("/home", wsMessenger.HomePage) // TODO: This can be removed later on...
	http.HandleFunc("/ws",  wsMessenger.WSEndpoint)
}

func main() {
	fmt.Println("Hello, World!")
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
