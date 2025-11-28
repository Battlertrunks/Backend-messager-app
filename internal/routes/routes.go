package routes

import (
	"net/http"

	"github.com/Battlertrunks/internal/app"
	"github.com/Battlertrunks/internal/wsMessenger"
	"github.com/gin-gonic/gin"
)

func Routes(app *app.Application) (*gin.Engine, error) {
	route := gin.Default()
	
	route.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Its OK",
		})
	})
	route.GET("/ws", func(ctx *gin.Context) {
		wsMessenger.WSEndpoint(ctx.Writer, ctx.Request)
	})

	route.GET("/api/v1/retrieve-user/:id", app.UserHandler.HandleRetrieveUser)
	route.POST("/api/v1/create-user", app.UserHandler.HandleCreateUser)

	return route, nil
}