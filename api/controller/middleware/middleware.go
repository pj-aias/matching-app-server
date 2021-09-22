package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/pj-aias/matching-app-server/auth"
	"github.com/pj-aias/matching-app-server/db"
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

		// validate signature
		// FIXME: extract to another midldeware
		signature := c.GetHeader("AnoMatch-Signature")
		gpk, err := db.GetUserGpk(uint(userId))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not find GPK"})
			return
		}

		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// body (Reader) will be consumed if once read, so reset the data
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		ok, err := auth.VerifySignature(body, signature, gpk)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to validate signature"})
			return
		}
		if !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid signature"})
			return
		}

		c.Next()
	}
}
