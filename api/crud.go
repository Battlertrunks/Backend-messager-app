package api

import (
	"database/sql"
	"net/http"

	"github.com/Battlertrunks/database"
	"github.com/Battlertrunks/wsMessenger"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userStore database.UserStore
}

func ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
			"message": "Its OK",
		})
}

func Crud(r *gin.Engine, db *sql.DB) {
	// Checks if the server is alive and responsive

	/* TODO:
		TODO:
			1. create new users
				- Need to research how to store passwords by hash
				- Maybe lightly look into MFA or SMS login too?
			2. login existing users
				- Research how to identify session, will require middleware.
				- JWT token or some other way?
			3. figure out websocket connection to be linked between certain clients (DM or group message)
			  - Need to save contact as friends
				- May need to make a friends list DB related to the user entry on the DB
				- How to find users to make as friends? (search on username or email for them?)
			4. Make a login/sign-up page to match the server's sign-up / login
	*/

	// Requests to establish a websocket connection
	r.GET("/ws", func(ctx *gin.Context) {
		wsMessenger.WSEndpoint(ctx.Writer, ctx.Request)

		ctx.JSON(http.StatusOK, gin.H{ "message": "CONNECTED WS" })
	})


	type NewUser struct {
		Username 	string `json:username`
		Email 	 	string `json:email`
		Password  string `json:password`
	}
	// Sign up / create user POST
	r.POST("/api/v1/create-user", func(ctx *gin.Context) {
		var newUser NewUser
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

		
		

		// We need the email, username, and password as a start... Fail if one of them are not included

	})
}

func (uh *UserHandler) HandleCreateUser(r *gin.Engine, ctx *gin.Context) {
	
}