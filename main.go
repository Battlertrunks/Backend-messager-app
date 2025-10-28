package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// TODO:
// Evaluate effort to adding a profiles for direct messaging and DB (Extra credit work)

type message struct {
	Message string `json:"message"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { return true },
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// upgrades connection to a websocket connection!
	log.Println("Started")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("err", err)
	}

	defer ws.Close()

	clients[ws] = true

	log.Println("Client Connected")
	// err = ws.WriteMessage(1, []byte("Hello, Client!"))
	if err != nil {
		log.Println("WriteMessage", err)
	}

	// log.Println("websocket", ws)
	reader(ws)
}

// Circle back on this to look for certain clients for direct messaging (extra credit work)
var clients = make(map[*websocket.Conn]bool) // Tracks active clients

func reader(conn *websocket.Conn) {
	// infinite loop to detect incoming messages
	for {
		// read in a message
		// TODO: Stop the message being read if it is from the original sender...
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("reader", err)
			return
		}

		// print out that message
		fmt.Println(string(msg))

		writer(msg)
	}
}

/*
TODO:
1. Make the write identify who sent the message to make a messenger and receiver on the client
*/
func writer(msg []byte) {
	// broadcast to all clients
	for client := range clients {
		if err := client.WriteMessage(websocket.TextMessage, msg); err != nil {
			fmt.Printf("Broadcast error:", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	fmt.Println("Hello, World!")
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
