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

type PasswordHash struct {
	ID   int
	Hash []byte `gorm:"size:72"`
	User User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func autoMigrate(database *gorm.DB) {
	database.AutoMigrate(&User{})
	database.AutoMigrate(&PasswordHash{})
}
