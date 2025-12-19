package middleware

import (
	"fmt"
	"log"
	"strings"

	"github.com/Battlertrunks/internal/app"
	"github.com/Battlertrunks/internal/store"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func SessionsMiddleware(app *app.Application) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Println(ctx.Request.URL.Path, strings.Contains(ctx.Request.URL.Path, "/api/v1/login"), ctx.Request.Method != "POST")
		if strings.Contains(ctx.Request.URL.Path, "/api/v1/login") && ctx.Request.Method == "POST" {
			ctx.Next()
			return
		}

		app.UserHandler.AuthorizeUser(ctx)

		ctx.Next()
	}
}
