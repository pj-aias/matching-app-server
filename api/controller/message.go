package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pj-aias/matching-app-server/db"
	"gorm.io/gorm"
)

type Message struct {
	Id         uint   `json:"id"`
	ChatroomId uint   `json:"chatroom_id"`
	User       User   `json:"user"`
	Content    string `json:"content"`
}

type MessageResponse struct {
	Message Message `json:"message"`
}

type MessagesResponse struct {
	Message []Message `json:"messages"`
}

func fromDBMessage(raw db.Message) Message {
	if raw.User == (db.User{}) {
		raw.User, _ = db.GetUser(uint64(raw.UserID))
	}

	return Message{
		Id:      raw.ID,
		User:    fromRawData(raw.User),
		Content: raw.Content,
	}
}

func fromDBMessages(rawMessages []db.Message) []Message {
	messages := make([]Message, len(rawMessages))

	for i, m := range rawMessages {
		messages[i] = fromDBMessage(m)
	}

	return messages
}

func AddMessage(c *gin.Context) {
	type postData struct {
		Content string
	}

	data := postData{}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, ok := c.MustGet("userId").(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, "invalid user id")
		return
	}

	roomId, err := strconv.ParseUint(c.Param("roomId"), 0, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if data.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty content is not allowed"})
		return
	}

	createdPost, err := db.CreateMessage(uint(userId), uint(roomId), data.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	post := fromDBMessage(createdPost)
	response := MessageResponse{post}
	c.JSON(http.StatusOK, response)
}

func ShowMessages(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 0, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rawMessage, err := db.GetMessage(uint(id))

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error})
		return
	}

	post := fromDBMessage(rawMessage)

	response := MessageResponse{post}
	c.JSON(http.StatusOK, response)
}

func ShowRooms(c *gin.Context) {
	type param struct {
		Count int
	}

	data := param{}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomId, err := strconv.ParseUint(c.Param("roomId"), 0, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	count := data.Count

	recentPosts, err := db.GetMessages(uint(roomId), count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	posts := MessagesResponse{fromDBMessages(recentPosts)}

	c.JSON(http.StatusOK, posts)
}

func UpdateMessageContent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 0, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	type param struct {
		Content string
	}

	data := param{}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, ok := c.MustGet("userId").(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, "invalid user id")
		return
	}

	old, err := db.GetMessage(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	if old.UserID != userId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you cannot update a post that was created by other users"})
		return
	}

	updatedPost, err := db.UpdateMessageContent(uint(id), data.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	post := fromDBMessage(updatedPost)
	response := MessageResponse{post}

	c.JSON(http.StatusOK, response)
}

func DeleteMessage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 0, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, ok := c.MustGet("userId").(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, "invalid user id")
		return
	}

	target, err := db.GetMessage(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	if target.UserID != userId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you cannot update a post that was created by other users"})
		return
	}

	err = db.DeleteMessage(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}