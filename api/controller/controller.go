package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pj-aias/matching-app-server/db"
)

func UserShow(c *gin.Context) {
	c.JSON(http.StatusOK, db.User{})
}
