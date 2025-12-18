package api

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Battlertrunks/internal/store"
	"github.com/gin-contrib/sessions"
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

	sessionToken := generateToken(32)
	tokenExpiresAt := time.Hour * 24

	userLoginResponse, err := uh.userStore.Login(userLogin, sessionToken, tokenExpiresAt)
	if err != nil {
		fmt.Println("500 ERROR")
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"message": "Unable to store session data"})
		return
	}

	fmt.Println("Pre session")

	session := sessions.Default(ctx)
	fmt.Println("Default")
	session.Set("userID", userLoginResponse.UserID)
	session.Set("username", userLoginResponse.Username)
	fmt.Println("Sets")
	err = session.Save()

	fmt.Println("RUN")

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save to session"})
		return
	}

	// Set session token
	ctx.SetCookie(
		"session_token",
		sessionToken,
		int(tokenExpiresAt),
		"/",
		"localhost", // use a env variable later on
		true,
		true,
	)

	ctx.JSON(http.StatusCreated, gin.H{"username": userLogin.Username, "password": userLogin.Password})
}

func generateToken(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}

	return base64.URLEncoding.EncodeToString(bytes)
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

	// if err := uh.AuthorizeUser(ctx, body.Username); err != nil {
	// 	ctx.AbortWithError(http.StatusUnauthorized, err)
	// 	return
	// }

	ctx.JSON(http.StatusCreated, gin.H{"status": "..."})
}

func (uh *UserHandler) GenerateCSRFToken(ctx *gin.Context) {
	csrfToken := generateToken(32)

	ctx.SetCookie(
		"csrf_token",
		csrfToken,
		60*60, // *60*24, // a day: 24 hours
		"/",
		"localhost", // use a env variable later on
		false,
		true,
	)

	ctx.JSON(http.StatusOK, gin.H{"message": "csrf token created", "csrf_token": csrfToken})

}

// TODO:
// - Create a new table "sessions" to store user_id, session_token, created_at, expires_at
// - Update the Middleware in internal/routes to look into the sessions based on user_id to compare the
//   session token with what is in the cookie
// Remove this func, AuthorizeUser and the one in internal/store
