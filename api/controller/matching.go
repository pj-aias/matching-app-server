package controller

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pj-aias/matching-app-server/db"
)

type MatchResponse struct {
	MatchedUser User     `json:"matched_user"`
	Chatroom    Chatroom `json:"chatroom"`
}

func MakeMatch(c *gin.Context) {
	userId, ok := c.MustGet("userId").(int)
	if !ok {
		e := fmt.Sprintf("invalid user id: %v", c.MustGet("userId"))
		c.JSON(http.StatusBadRequest, gin.H{"error": e})
		return
	}

	targetUsers, err := getMatchUsers(uint(userId))
	if err != nil {
		e := fmt.Sprintf("failed to cummunicate with the database: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": e})
		return
	}

	if len(targetUsers) <= 0 {
		c.JSON(http.StatusOK, gin.H{"error": "no user matched"})
		return
	}

	matchedUser := selectUser(targetUsers)

	rawChatroom, err := db.CreateRoom([]uint{uint(userId), matchedUser.ID})
	if err != nil {
		e := fmt.Sprintf("failed to cummunicate with the database: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": e})
		return
	}
	chatroom := fromDBRoom(rawChatroom)

	response := MatchResponse{matchedUser, chatroom}

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
