package main

import (
	"github.com/gin-gonic/gin"

	"github.com/pj-aias/matching-app-server/controller"
)

func main() {
	r := gin.Default()

	r.GET("user/:id", controller.UserShow)
	r.POST("user", controller.UserAdd)
	r.PATCH("user", controller.UserUpdate)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
