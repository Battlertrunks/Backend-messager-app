package routes

import (
	"net/http"

	"github.com/Battlertrunks/internal/app"
	"github.com/Battlertrunks/internal/middleware"
	"github.com/Battlertrunks/internal/wsMessenger"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func Routes(app *app.Application) (*gin.Engine, error) {
	route := gin.Default()

	store := cookie.NewStore([]byte("the-secret-key"))
	route.Use(sessions.Sessions("appsession", store))
	route.Use(middleware.CSRFMiddleware())
	route.Use(middleware.SessionsMiddleware(app))

	// TODO: Make the category of routes to be in their own go file that meets here with the Routes
	// This can avoid too much clutter in the routes

	route.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Its OK",
		})
	})

	route.GET("/ws", func(ctx *gin.Context) {
		wsMessenger.WSEndpoint(ctx.Writer, ctx.Request)
	})

	// --- CSRF token creation ---
	route.GET("/api/v1/generate-csrf-token", app.UserHandler.GenerateCSRFToken)

	// --- User Routes ---
	route.GET("/api/v1/retrieve-user/:id", app.UserHandler.HandleRetrieveUser)
	route.POST("/api/v1/create-user", app.UserHandler.HandleCreateUser)
	route.POST("/api/v1/login", app.UserHandler.HandleUserLogin)
	route.POST("/api/v1/logout", app.UserHandler.HandleUserLogout)
	route.POST("/api/v1/protected") // Do I need that?

	return route, nil
}
