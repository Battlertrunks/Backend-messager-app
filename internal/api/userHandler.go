package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

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
	fmt.Println("TEST CREATE")
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

func (uh *UserHandler) HandleRetrieveUser(ctx *gin.Context) {
	fmt.Println("TEST")
	userIdStr, exists := ctx.Params.Get("id")
	// If the param is not existent in the params
	fmt.Println(exists, userIdStr)
	if !exists {
		ctx.AbortWithStatus(404)
	}

	fmt.Println(userIdStr)

	// Convert the param from string to uint64
	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		ctx.AbortWithStatus(500)
	}

	user, err := uh.userStore.GetUser(userId)
	if err != nil {
		ctx.AbortWithStatus(500)
	}

	fmt.Println(user)

	// If user is not found in the database
	if user.ID == 0 {
		ctx.AbortWithStatus(404)
	}

	ctx.JSON(http.StatusAccepted, user)
}