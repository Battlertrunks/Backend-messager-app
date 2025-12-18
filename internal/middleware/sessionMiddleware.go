package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Battlertrunks/internal/store"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func (uh *UserHandler) SessionsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !strings.Contains(ctx.Request.URL.RawPath, "/api/v1/login") && ctx.Request.Method != "POST" {
			ctx.Next()
			return
		}

		SessionToken, err := ctx.Cookie("session_token")
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			ctx.Next()
			return
		}

		session := sessions.Default(ctx)
		sessionUserID := session.Get("userID")

		UserID := sessionUserID.(int)

		fmt.Printf("Create the cookie evaluation to the store for session token: %v\n", SessionToken)
		err = uh.userStore.AuthorizeUser(store.AuthorizeData{UserID: UserID, SessionToken: SessionToken})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			ctx.Next()
			return
		}

		ctx.Next()
	}
}
