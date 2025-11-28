package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Battlertrunks/internal/store"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userStore store.UserStore
	logger *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger: 	 logger,
	}
}

type NewUser struct {
	Username string `json:"username"`
	Email 	 string `json:"email"`
	Password string `json:"password"`
}

func (uh *UserHandler) HandleCreateUser(ctx *gin.Context) {
	var newUser store.User
		if err := ctx.ShouldBindJSON(&newUser); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{ "message": "Error parsing JSON" })
			return
		}
		
		if newUser.Email == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{ "message": "Missing 'email' value..." })
			return
		}

		if newUser.Username == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{ "message": "Missing 'username' value..." })
			return
		}

		if newUser.Password == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{ "message": "Missing 'password' value..." })
			return
		}

		createdUser, err := uh.userStore.CreateUser(&newUser)
		if err != nil {
			fmt.Println(err)
			ctx.JSON(http.StatusBadRequest, gin.H{ "message": "Could not create user" })
			return
		}

		ctx.JSON(http.StatusCreated, createdUser)
}