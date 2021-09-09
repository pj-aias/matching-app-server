package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pj-aias/matching-app-server/auth"
	"github.com/pj-aias/matching-app-server/db"
	"gorm.io/gorm"
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

func fromDBUsers(rawUsers []db.User) []User {
	users := make([]User, len(rawUsers))
	for i, u := range rawUsers {
		users[i] = fromRawData(u)
	}
	return users
}

func UserShow(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 0, 64)
	if err != nil {
		e := fmt.Sprintf("invalid user id (%v): %v", c.Param("id"), err)
		c.JSON(http.StatusBadRequest, gin.H{"error": e})
		return
	}

	result, err := db.GetUser(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		e := fmt.Sprintf("user with that id was not found: %v", id)
		c.JSON(http.StatusNotFound, gin.H{"error": e})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "an error occured while getting user from database: " + err.Error()})
		return
	}

	user := fromRawData(result)
	c.JSON(http.StatusOK, user)
}

func UserAdd(c *gin.Context) {
	type postData struct {
		Username  string `binding:"required"`
		Password  string `binding:"required"`
		Signature string `binding:"required"`
	}
	type responseData struct {
		User  User   `json:"user"`
		Token string `json:"token"`
	}

	data := postData{}

	// todo empty signature passes
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get post data: " + err.Error()})
		return
	}

	_, err := db.LookupUser(data.Username)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "a user with that username already exists"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// normal error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "an error occured while contacting to database: " + err.Error()})
		return
	}

	// validate password and generate hash before inserting user data into DB
	passwordHash, err := auth.GeneratePasswordHash(data.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to generate password hash: " + err.Error()})
		return
	}

	userParam := db.User{
		Username: data.Username,
	}
	createdUser, err := db.AddUser(userParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user: " + err.Error()})
		return
	}
	user := fromRawData(createdUser)

	_, err = db.AddPasswordHash(uint64(user.ID), passwordHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to saving hashed password to database: " + err.Error()})
		return
	}

	token, err := auth.CreateToken(int(user.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to issue a auth credential: " + err.Error()})
		return
	}

	response := responseData{
		User:  user,
		Token: token,
	}

	c.JSON(http.StatusOK, response)
}

func UserUpdate(c *gin.Context) {
	type updateData struct {
		Avatar string `json:",omitempty"`
		Bio    string `json:",omitempty"`
	}

	data := updateData{}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get post data: " + err.Error()})
		return
	}

	userId, ok := c.MustGet("userId").(int)
	if !ok {
		e := fmt.Sprintf("invalid user id: %v", c.MustGet("userId"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": e})
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
		User  User   `json:"user"`
		Token string `json:"token"`
	}

	data := postData{}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get post data: " + err.Error()})
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

	err = auth.ValidatePassword(hash.Hash, data.Password)
	if errors.Is(err, &auth.ErrPasswordDidNotMatch) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "password did not match"})
	} else if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := auth.CreateToken(int(user.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := responseData{
		User:  fromRawData(user),
		Token: token,
	}

	c.JSON(http.StatusOK, response)
}
