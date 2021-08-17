package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pj-aias/matching-app-server/db"
)

type Follow struct {
	Target    uint
	Followed  bool
	Following bool
}

func FollowUser(c *gin.Context) {
	destUserId, err := strconv.ParseUint(c.Param("id"), 0, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	srcUserId, ok := c.MustGet("userId").(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, "invalid user id")
		return
	}

	_, err = db.CreateFollow(uint(srcUserId), uint(destUserId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	response := Follow{}
	response.Target = uint(destUserId)
	response.Following = true
	response.Followed, err = db.DoesFollow(uint(destUserId), uint(srcUserId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}
