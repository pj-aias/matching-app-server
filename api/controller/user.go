package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pj-aias/matching-app-server/auth"
	"github.com/pj-aias/matching-app-server/db"
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
	type responseData struct {
		User User `json:"user"`
		Token string `json:"token"`
	}

	data := postData{}

	// todo empty signature passes
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	oldUser, _ := db.LookupUser(data.Username)
	if oldUser != (db.User{}) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "a user with that username already exists"})
		return
	}

	// validate password and generate hash before inserting user data into DB
	passwordHash, err := auth.GeneratePasswordHash(data.Password)
	if err != nil {
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

	_, err = db.AddPasswordHash(uint64(user.ID), passwordHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := auth.CreateToken(int(user.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	response := responseData {
		User: user,
		Token: token,
	}

	c.JSON(http.StatusOK, response)
}

func UserUpdate(c *gin.Context) {
	type updateData struct {
		Avatar    string `json:",omitempty"`
		Bio       string `json:",omitempty"`
	}

	data := updateData{}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, ok := c.MustGet("userId").(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, "invalid user id")
		return
	}

	userData := db.User{}
	userData.Avatar = data.Avatar
	userData.Bio = data.Bio

	_, err := db.UpdateUser(uint(userId), userData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	user, err := db.GetUser(uint64(userId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, fromRawData(user))
}

func Login(c *gin.Context) {
	type postData struct {
		Username  string
		Password  string
		Signature string
	}
	type responseData struct {
		User User `json:"user"`
		Token string `json:"token"`
	}

	data := postData{}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := db.LookupUser(data.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := db.GetPasswordHash(uint64(user.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := auth.ValidatePassword(hash.Hash, data.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := auth.CreateToken(int(user.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := responseData {
		User: fromRawData(user),
		Token: token,
	}

	c.JSON(http.StatusOK, response)
}