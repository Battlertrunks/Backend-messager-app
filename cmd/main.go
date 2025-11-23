package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Battlertrunks/internal/app"
	"github.com/Battlertrunks/internal/routes"
	"github.com/gin-gonic/gin"
)

func setupRoutes(app *app.Application) (*gin.Engine)  {
	// Initialize the REST API
	r, err := routes.Routes(app)
	if err != nil {
		log.Fatal("Database Failed")
		panic(err)
	}

	// Initialize the WebSocket connection
	// TODO: We should start the websocket based on when the user is in a live chat connection with
	// a another receiver
	// http.HandleFunc("/home", wsMessenger.HomePage) // TODO: This can be removed later on...
	// http.HandleFunc("/ws",  wsMessenger.WSEndpoint)

	return r
}

func main() {
	
	port := 8080

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	routes := setupRoutes(app)

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: routes,
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
