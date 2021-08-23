package controller

import (
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pj-aias/matching-app-server/db"
)

type MatchResponse struct {
	User User `json:"user"`
}

func MakeMatch(c *gin.Context) {
	targetUsers, err := getMatchUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(targetUsers) <= 0 {
		c.JSON(http.StatusOK, gin.H{"error": "no user matched"})
		return
	}

	matchedUser := selectUser(targetUsers)
	response := MatchResponse{matchedUser}

	c.JSON(http.StatusOK, response)
}

func getMatchUsers() ([]User, error) {
	rawUsers, err := db.GetAllUsers()
	if err != nil {
		return nil, err
	}

	users := fromDBUsers(rawUsers)
	return users, nil
}

func selectUser(targets []User) User {
	return randomPick(targets)
}

func randomPick(targets []User) User {
	i := rand.Intn(len(targets))
	return targets[i]
}
