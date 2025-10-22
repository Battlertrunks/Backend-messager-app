package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

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

	log.Println("Client Connected")
	err = ws.WriteMessage(1, []byte("Hello, Client!"))
	if err != nil {
		log.Println("WriteMessage", err)
	}

	// log.Println("websocket", ws)
	reader(ws)
}

func reader(conn *websocket.Conn) {
	// infinite loop to detect incoming messages
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("reader", err)
			return
		}

		// print out that message
		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println("reader write", err)
			return
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
