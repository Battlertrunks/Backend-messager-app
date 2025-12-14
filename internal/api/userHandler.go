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
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

type NewUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Username string
	Email    string
}

func (uh *UserHandler) HandleCreateUser(ctx *gin.Context) {
	var newUser store.User
	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Error parsing JSON"})
		return
	}

	// TODO: Can these checks be more consolidated or universal?
	if newUser.Email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Missing 'email' value..."})
		return
	}

	if newUser.Username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Missing 'username' value..."})
		return
	}

	if newUser.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Missing 'password' value..."})
		return
	}

	if len(newUser.Username) < 8 || len(newUser.Password) < 8 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Username and Password cannot be less than 8 characters"})
		return
	}
	// -----------------------------------------------------------

	createdUser, err := uh.userStore.CreateUser(&newUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdUser)
}

func (uh *UserHandler) HandleRetrieveUser(ctx *gin.Context) {
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

func (uh *UserHandler) HandleUserLogin(ctx *gin.Context) {
	var userLogin store.LoginData
	if err := ctx.ShouldBindJSON(&userLogin); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	sessionToken, err := uh.userStore.Login(userLogin)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.SetCookie(
		"session_cookie",
		sessionToken,
		60*60*24, // a day: 24 hours
		"/",
		"localhost", // use a env variable later on
		true,
		true,
	)
}

func (uh *UserHandler) HandleUserLogout(ctx *gin.Context) {
	// TODO: Write the user logout functionality from cookies

}
