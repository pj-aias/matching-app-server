package db

import "gorm.io/gorm"

func GetUser(id uint64) (User, error) {
	user := User{
		gorm.Model{ID: uint(id)},
	}
	database.Take(&user)

	return user, database.Error
}
