package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pj-aias/matching-app-server/auth"
)

func AuthorizeToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer"
		authHeader := c.GetHeader("Authorization")
		tokenString := authHeader[len(BEARER_SCHEMA):]
		userId, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Set("userId", userId)
		c.Next()
	}
}