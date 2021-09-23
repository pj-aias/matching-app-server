package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/pj-aias/matching-app-server/controller"
	"github.com/pj-aias/matching-app-server/controller/middleware"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})
	r.POST("user", controller.UserAdd)

	// require existing users to send signature
	existingUsers := r.Group("/")
	existingUsers.Use(middleware.VerifySignature)

	{
		existingUsers.POST("login", controller.Login)

		// Authorization with JWT required
		authRequired := existingUsers.Group("/")
		authRequired.Use(middleware.AuthorizeToken())

		{
			authRequired.GET("user/:id", controller.UserShow)
			authRequired.PATCH("user", controller.UserUpdate)

			authRequired.GET("follow/:id", controller.ShowFollow)
			authRequired.POST("follow/:id", controller.FollowUser)
			authRequired.DELETE("follow/:id", controller.UnfollowUser)
			authRequired.GET("followers", controller.ShowFollowers)
			authRequired.GET("followees", controller.ShowFollowees)

			authRequired.POST("message", controller.CreateRoom)
			authRequired.GET("message/rooms", controller.ShowRooms)
			authRequired.POST("message/:roomId", controller.AddMessage)
			authRequired.GET("message/:roomId", controller.ShowMessages)

			authRequired.POST("matching", controller.MakeMatch)

		}
	}

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
