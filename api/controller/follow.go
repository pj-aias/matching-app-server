package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pj-aias/matching-app-server/db"
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
		e := fmt.Sprintf("invalid user id: %v", c.MustGet("userId"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": e})
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
		e := fmt.Sprintf("invalid user id: %v", c.MustGet("userId"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": e})
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
		e := fmt.Sprintf("invalid user id: %v", c.MustGet("userId"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": e})
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

type Followees struct {
	Users []User `json:"followees"`
}

func ShowFollowees(c *gin.Context) {
	source, ok := c.MustGet("userId").(int)
	if !ok {
		e := fmt.Sprintf("invalid user id: %v", c.MustGet("userId"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": e})
		return
	}

	followings, err := db.GetFollowing(uint(source))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	followees, err := follows2users(followings, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	response := Followees{
		Users: followees,
	}

	c.JSON(http.StatusOK, response)
}

type Followers struct {
	Users []User `json:"followers"`
}

func ShowFollowers(c *gin.Context) {
	source, ok := c.MustGet("userId").(int)
	if !ok {
		e := fmt.Sprintf("invalid user id: %v", c.MustGet("userId"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": e})
		return
	}

	followeds, err := db.GetFollowed(uint(source))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	followers, err := follows2users(followeds, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	response := Followers{
		Users: followers,
	}

	c.JSON(http.StatusOK, response)
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
		Target:    target,
		Following: following,
		Followed:  followed,
	}, nil
}

func follows2users(follows []db.Follow, forDest bool) ([]User, error) {
	usersId := make([]uint, len(follows))
	if forDest {
		for i, f := range follows {
			usersId[i] = uint(f.DestUserID)
		}
	} else {
		for i, f := range follows {
			usersId[i] = uint(f.SourceUserID)
		}
	}

	dbUsers, err := db.GetUsers(usersId)
	if err != nil {
		return nil, err
	}

	users := make([]User, len(dbUsers))
	for i, u := range dbUsers {
		users[i] = fromRawData(u)
	}

	return users, nil
}
