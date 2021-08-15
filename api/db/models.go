package db

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name   string `gorm:"size:255"`
	Avatar string `gorm:"size:255"`
	Bio    string `gorm:"size:255"`
}

func autoMigrate(database *gorm.DB) {
	database.AutoMigrate(&User{})
}