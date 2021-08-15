package db

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"size:32"`
	Avatar   string `gorm:"size:255"`
	Bio      string `gorm:"size:255"`
}

func autoMigrate(database *gorm.DB) {
	database.AutoMigrate(&User{})
}
