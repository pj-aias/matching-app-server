package main

import (
	"github.com/gin-gonic/gin"

	"github.com/pj-aias/matching-app-server/controller"
	"github.com/pj-aias/matching-app-server/controller/middleware"
)

func main() {
	r := gin.Default()
	r.POST("user", controller.UserAdd)
	r.POST("login", controller.Login)

	authRequired := r.Group("/")
	authRequired.Use(middleware.AuthorizeToken())

	{
		authRequired.GET("user/:id", controller.UserShow)
		authRequired.PATCH("user", controller.UserUpdate)

		authRequired.GET("follow/:id", controller.ShowFollow)
		authRequired.POST("follow/:id", controller.FollowUser)
	}

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
