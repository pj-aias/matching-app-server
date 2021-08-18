package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pj-aias/matching-app-server/db"
	"gorm.io/gorm"
)

type Follow struct {
	Target    uint `json:"target"`
	Followed  bool `json:"followed"`
	Following bool `json:"following"`
}

func ShowFollow(c *gin.Context) {
	target, err := strconv.ParseUint(c.Param("id"), 0, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	srcUserId, ok := c.MustGet("userId").(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, "invalid user id")
		return
	}

	response, err := getFollowFromDB(uint(srcUserId), uint(target))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
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

	response, err := getFollowFromDB(uint(srcUserId), uint(destUserId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

func UnfollowUser(c *gin.Context) {
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

	err = db.DestroyFollow(uint(srcUserId), uint(destUserId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	response, err := getFollowFromDB(uint(srcUserId), uint(destUserId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

type Follwees struct {
	Users []User `json:"followees"`
}

func ShowFollowees(c *gin.Context) {
	source, ok := c.MustGet("userId").(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, "invalid user id")
		return
	}

	followings, err := db.GetFollowing(source)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	followees := 
}

func getFollowFromDB(source, target uint) (Follow, error) {
	following, err := db.DoesFollow(source, target)
	if err != nil {
		return Follow{}, err
	}

	followed, err := db.DoesFollow(target, source)
	if err != nil {
		return Follow{}, err
	}

	return Follow{
		Target: target,
		Following: following,
		Followed: followed,
	}, nil
}

func follows2users(follows []Follow) ([]User, err) {
	count := len(follows)
	usersId := make([]uint, count)
	for i, f := range follows {
		usersId[i] = f.Target
	}

	dbUsers, err := db.GetUsers(usersId)
	if err != nil {
		return nil, err
	}

	users := make([]User, count)
	for i, u := range dbUsers {
		users[i] = fromRawData(u)
	}

	return users, nil
}