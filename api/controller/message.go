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
	Id    uint   `json:"id"`
	Users []User `json:"users"`
}

type ChatroomResponse struct {
	Chatroom Chatroom `json:"chatroom"`
}

type ChatroomsResponse struct {
	Chatrooms []Chatroom `json:"chatrooms"`
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
	users := fromDBUsers(rawRoom.Users)

	return Chatroom{
		Id:    rawRoom.ID,
		Users: users,
	}
}

func fromDBRooms(rawRooms []db.Chatroom) []Chatroom {
	rooms := make([]Chatroom, len(rawRooms))
	for i, r := range rawRooms {
		rooms[i] = fromDBRoom(r)
	}

	return rooms
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

	room, err := db.GetRoom(uint(roomId))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !room.ContainsUser(uint(userId)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you cannot post a message to room you are not in"})
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
	userId, ok := c.MustGet("userId").(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id"})
		return
	}

	rawRooms, err := db.GetRooms(uint(userId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rooms := fromDBRooms(rawRooms)
	response := ChatroomsResponse{rooms}
	c.JSON(http.StatusOK, response)
}

func UpdateMessageContent(c *gin.Context) {
	type param struct {
		TargetMessageId int
		Content         string
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

	old, err := db.GetMessage(uint(data.TargetMessageId))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	if old.UserID != userId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you cannot update a post that was created by other users"})
		return
	}

	_, err = db.UpdateMessageContent(uint(data.TargetMessageId), data.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rawMessage, err := db.GetMessage(uint(data.TargetMessageId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	msg := fromDBMessage(rawMessage)
	response := MessageResponse{msg}

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
