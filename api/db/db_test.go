package db_test

import (
	"github.com/pj-aias/matching-app-server/db"
)

func sampleUser(id uint, username string) db.User {
	user := db.User{}
	user.ID = id
	user.Username = username

	return user
}
