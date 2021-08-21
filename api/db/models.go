package db

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"size:32,unique"`
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

type Message struct {
	gorm.Model
	ChatroomId int
	Chatroom   Chatroom `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UserID     int
	User       User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Content    string
}

type Chatroom struct {
	gorm.Model
	Users    []User `gorm:"many2many:chatroom_users;"`
	Messages []Message
}

func autoMigrate(database *gorm.DB) {
	database.AutoMigrate(&User{})
	database.AutoMigrate(&PasswordHash{})
	database.AutoMigrate(&Follow{})
	database.AutoMigrate(&Message{})
	database.AutoMigrate(&Chatroom{})
}
