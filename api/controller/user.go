package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pj-aias/matching-app-server/db"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
}

func fromRawData(raw db.User) User {
	return User{
		ID:       raw.ID,
		Username: raw.Username,
		Avatar:   raw.Avatar,
		Bio:      raw.Bio,
	}
}

func UserShow(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 0, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.GetUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	user := fromRawData(result)
	c.JSON(http.StatusOK, user)
}

func UserAdd(c *gin.Context) {
	type postData struct {
		Username  string
		Password  string
		Signature string
	}

	data := postData{}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userParam := db.User{
		Username: data.Username,
	}
	createdUser, err := db.AddUser(userParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user := fromRawData(createdUser)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err = db.AddPasswordHash(uint64(user.ID), hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, user)
}

func UserUpdate(c *gin.Context) {
	type updateData struct {
		Name      *string
		Avatar    *string
		Bio       *string
		Signature string
	}

	data := updateData{}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}
