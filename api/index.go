package api

import (
	"github.com/gin-gonic/gin"
)

func RunGin() {
	r := gin.Default()

	Crud(r)

	r.Run(":8080")
}