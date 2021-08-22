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

type Chatroom struct {
	Id    uint  `json:"id"`
	Users []int `json:"users"`
}

type ChatroomResponse struct {
	Chatroom Chatroom `json:"chatroom"`
}

func fromDBMessage(raw db.Message) Message {
	if raw.User == (db.User{}) {
		raw.User, _ = db.GetUser(uint64(raw.UserID))
	}

	return Message{
		Id:         raw.ID,
		User:       fromRawData(raw.User),
		Content:    raw.Content,
		ChatroomId: uint(raw.ChatroomId),
	}
}

func fromDBMessages(rawMessages []db.Message) []Message {
	messages := make([]Message, len(rawMessages))

	for i, m := range rawMessages {
		messages[i] = fromDBMessage(m)
	}

	return messages
}

func fromDBRoom(rawRoom db.Chatroom) Chatroom {
	userIds := make([]int, len(rawRoom.Users))
	for i, u := range rawRoom.Users {
		userIds[i] = int(u.ID)
	}

	return Chatroom{
		Id:    rawRoom.ID,
		Users: userIds,
	}
}

func CreateRoom(c *gin.Context) {
	type postData struct {
		Target int
	}

	data := postData{}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	targetId := data.Target
	userId, ok := c.MustGet("userId").(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, "invalid user id")
		return
	}

	createdRoom, err := db.CreateRoom([]uint{uint(userId), uint(targetId)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	room := fromDBRoom(createdRoom)
	response := ChatroomResponse{room}
	c.JSON(http.StatusOK, response)
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

	_, err = db.GetRoom(uint(roomId))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	/*
		TODO: Block request from users who is not in chatroom members.
		      Currently `ContainsUser` not working

		if !room.ContainsUser(uint(userId)) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "you cannot post a message to room you are not in"})
			return
		}*/

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
	chatroomId, err := strconv.ParseUint(c.Param("roomId"), 0, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rawMessages, err := db.GetMessages(uint(chatroomId))

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "chatroom not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error})
		return
	}

	messages := fromDBMessages(rawMessages)

	response := MessagesResponse{messages}
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

	recentPosts, err := db.GetMessages(uint(roomId))
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
