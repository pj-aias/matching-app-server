package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

type Gm = string
type Gpk = string

type GpkRegistry interface {
	GetGpk(gm Gm) (Gpk, error)
}

func CombineGpk(gpks []Gpk) (Gpk, error) {
	// TODO implement CpmbineGpk process
	return gpks[0], nil
}

func GetGpkFromGms(registry GpkRegistry, gms []Gm) (Gpk, error) {
	gpks := []Gpk{}
	for _, gm := range gms {
		gpk, err := registry.GetGpk(gm)
		if err != nil {
			return "", fmt.Errorf("could not get gpk from %v: %v", gm, err)
		}

		gpks = append(gpks, gpk)
	}

	combinedGpk, err := CombineGpk(gpks)
	if err != nil {
		return "", fmt.Errorf("could not combine gpks: %v", err)
	}

	return combinedGpk, nil
}

func VerifySignature(gpkRegistry GpkRegistry) gin.HandlerFunc {
	// validate signature

	type gmsData = []string

	return func(c *gin.Context) {
		signature := c.GetHeader("AnoMatch-Signature")

		gmsJson := c.GetHeader("AIAS-GMs")
		gms := gmsData{}
		err := json.Unmarshal([]byte(gmsJson), &gms)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid format of gms"})
			return
		}

		gpk, err := GetGpkFromGms(gpkRegistry, gms)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "failed to get gpk: " + err.Error()})
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
