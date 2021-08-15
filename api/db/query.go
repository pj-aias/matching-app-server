package db

import "gorm.io/gorm"

func GetUser(id uint64) (User, error) {
	user := User{
		gorm.Model{ID: uint(id)},
	}
	result := database.Take(&user)

	return user, result.Error
}

func AddUser(user User) (User, error) {
	result := database.Create(&user)
	return user, result.Error
}
