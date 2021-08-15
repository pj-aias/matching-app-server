package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pj-aias/matching-app-server/db"
)

func UserShow(c *gin.Context) {
	c.JSON(http.StatusOK, db.User{})
}

func UserAdd(c *gin.Context) {
	type postData struct {
		Name      string
		Password  string
		Signature string
	}

	data := postData{}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
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
