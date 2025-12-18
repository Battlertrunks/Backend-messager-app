package middleware

import (
	"crypto/subtle"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CSRFMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method == "GET" {
			ctx.Next()
			return
		}

		csrfToken, err := ctx.Cookie("csrf_token")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Invalid session from user"})
			return
		}

		headerToken := ctx.GetHeader("X-CSRF-Token")
		if headerToken == "" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Invalid session from user"})
			return
		}

		fmt.Println("TOKENS", csrfToken, headerToken)
		if subtle.ConstantTimeCompare([]byte(csrfToken), []byte(headerToken)) != 1 {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "CSRF Token does not match"})
			return
		}

		ctx.Next()
	}
}
