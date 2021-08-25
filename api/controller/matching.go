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
	userId, ok := c.MustGet("userId").(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, "invalid user id")
		return
	}

	targetUsers, err := getMatchUsers(uint(userId))
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

func getMatchUsers(fromUserId uint) ([]User, error) {
	rawUsers, err := db.GetAllUsers()
	if err != nil {
		return nil, err
	}

	users := fromDBUsers(rawUsers)

	lastIdx := len(users) - 1
	// find source user
	for i, u := range users {
		if u.ID == fromUserId {
			users[i] = users[lastIdx]
			users = users[:lastIdx]
		}
	}

	return users, nil
}

func selectUser(targets []User) User {
	return randomPick(targets)
}

func randomPick(targets []User) User {
	i := rand.Intn(len(targets))
	return targets[i]
}
