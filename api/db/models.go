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
	Hash   []byte `gorm:"size:72"`
	UserID int
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Follow struct {
	gorm.Model
	SourceUserID int
	SourceUser   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	DestUserID   int
	DestUser     User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func autoMigrate(database *gorm.DB) {
	database.AutoMigrate(&User{})
	database.AutoMigrate(&PasswordHash{})
	database.AutoMigrate(&Follow{})
}
