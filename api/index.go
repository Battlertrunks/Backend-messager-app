package api

import (
	"github.com/Battlertrunks/database"
	"github.com/gin-gonic/gin"
)

func RunGin() error {
	sqlDB, err := database.Open()
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	r.GET("/ping", ping)

	Crud(r, sqlDB)

	r.Run(":8080")

	return nil
}