package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pj-aias/matching-app-server/auth"
)

func AuthorizeToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer"
		authHeader := c.GetHeader("Authorization")

		if !strings.HasPrefix(authHeader, BEARER_SCHEMA) {
			fmt.Fprintf(os.Stderr, "invalid Authorization header: '%v'\n", authHeader)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimSpace(authHeader[len(BEARER_SCHEMA):])
		userId, err := auth.ValidateToken(tokenString)
		if err != nil {
			fmt.Fprintf(os.Stderr, "authorization failed: '%v'\n", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("userId", userId)
		c.Next()
	}
}