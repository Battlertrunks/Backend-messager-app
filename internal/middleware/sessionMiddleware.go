package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func SessionsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !strings.Contains(ctx.Request.URL.RawPath, "/api/v1/login") && ctx.Request.Method != "POST" {
			ctx.Next()
			return
		}

		sessionToken, err := ctx.Cookie("session_token")
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			ctx.Next()
			return
		}

		// TODO call the user store to authorize the session token...
		fmt.Printf("Create the cookie evaluation to the store for session token: %v\n", sessionToken)

		// TODO When the call from the store is falsey err, run this.
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			ctx.Next()
			return
		}

		ctx.Next()
	}
}
