package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pj-aias/matching-app-server/auth"
)

func AuthorizeToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer"
		authHeader := c.GetHeader("Authorization")

		if !strings.HasPrefix(authHeader, BEARER_SCHEMA) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[len(BEARER_SCHEMA):]
		userId, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("userId", userId)
		c.Next()
	}
}