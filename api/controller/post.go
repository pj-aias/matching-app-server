package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pj-aias/matching-app-server/db"
	"gorm.io/gorm"
)

type Post struct {
	Id      uint   `json:"id"`
	User    User   `json:"user"`
	Content string `json:"content"`
}

type PostResponse struct {
	Post Post `json:"post"`
}

func fromDBPost(raw db.Post) Post {
	if raw.User == (db.User{}) {
		raw.User, _ = db.GetUser(uint64(raw.UserID))
	}

	return Post{
		Id:      raw.ID,
		User:    fromRawData(raw.User),
		Content: raw.Content,
	}
}

func PostAdd(c *gin.Context) {
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

	if data.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty content is not allowed"})
		return
	}

	createdPost, err := db.CreatePost(uint(userId), data.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	post := fromDBPost(createdPost)
	response := PostResponse{post}
	c.JSON(http.StatusOK, response)
}

func ShowPost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 0, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rawPost, err := db.GetPost(uint(id))

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error})
		return
	}

	post := fromDBPost(rawPost)

	response := PostResponse{post}
	c.JSON(http.StatusOK, response)
}
