package api

import (
	"net/http"

	"github.com/Battlertrunks/wsMessenger"
	"github.com/gin-gonic/gin"
)

func Crud(r *gin.Engine) {
	// Checks if the server is alive and responsive
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Its OK",
		})
	})

	// Requests to establish a websocket connection
	r.GET("/ws", func(ctx *gin.Context) {
		wsMessenger.WSEndpoint(ctx.Writer, ctx.Request)

		ctx.JSON(http.StatusOK, gin.H{ "message": "CONNECTED WS" })
	})
}