package main

import (
	"github.com/gin-gonic/gin"

	"github.com/pj-aias/matching-app-server/controller"
	"github.com/pj-aias/matching-app-server/db"
)

func main() {
	db.TestInsert(db.User{Name: "hoge"})

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("user", controller.UserShow)
	r.POST("user", controller.UserAdd)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
