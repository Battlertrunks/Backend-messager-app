package api

import (
	"errors"
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
	var userLogin store.UserData
	if err := ctx.ShouldBindJSON(&userLogin); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	sessionToken, csrfToken, err := uh.userStore.Login(userLogin)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Set session token
	ctx.SetCookie(
		"session_token",
		sessionToken,
		60, // *60*24, // a day: 24 hours
		"/",
		"localhost", // use a env variable later on
		true,
		true,
	)

	// Set cross site request forgery token
	// TODO: Find a middleware to wrap the CSRF with Gorrila
	ctx.SetCookie(
		"csrf_token",
		csrfToken,
		60, // *60*24, // a day: 24 hours
		"/",
		"localhost", // use a env variable later on
		true,
		false,
	)

	ctx.JSON(http.StatusCreated, gin.H{"username": userLogin.Username, "password": userLogin.Password})
}

func (uh *UserHandler) HandleUserLogout(ctx *gin.Context) {
	var body struct {
		Username string `json:"username"`
	}
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err := uh.AuthorizeUser(ctx, body.Username); err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "..."})
}

func (uh *UserHandler) AuthorizeUser(ctx *gin.Context, username string) error {
	var userAuth store.AuthorizeData
	userAuth.Username = username
	if userAuth.Username == "" {
		return errors.New("Unathorized access of user 1") // remove digit
	}

	var err error
	userAuth.SessionToken, err = ctx.Cookie("session_token")
	// TODO: Find a middleware to wrap the CSRF with Gorrila
	userAuth.CSRFToken = ctx.GetHeader("X-CSRF-Token")
	fmt.Printf("SessionToken: %v \n csrfToken: %v \n", userAuth.SessionToken, userAuth.CSRFToken)
	if userAuth.CSRFToken == "" || err != nil {
		return errors.New("Unathorized access of user 2") // remove digit
	}

	err = uh.userStore.AuthorizeUser(userAuth)
	if err != nil {
		return err
	}

	return nil // User is correctly authenticated
}
