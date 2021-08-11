package db

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	database *gorm.DB
)

func init() {
	user := "root"
	password := getPassword()
	host := "db"
	port := "3306"
	dbName := "service"
	config := gorm.Config{}

	database = connectDB(user, password, host, port, dbName, config)
}

func getPassword() string {
	filename := "/run/secrets/mysql_password"
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal("failed to open MySQL Password:", err)
	}

	defer f.Close()

	buf := make([]byte, 128)
	n, err := f.Read(buf)

	if n == 0 {
		log.Fatal("MySQL Password file empty")
	}
	if err != nil {
		log.Fatal("failed to read MySQL Password: ", err)
	}

	return string(buf)
}

func connectDB(user, password, host, port, dbName string, config gorm.Config) *gorm.DB {
	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbName + "?parseTime=true&loc=Asia%2FTokyo"
	db, err := gorm.Open(mysql.Open(dsn), &config)

	if err != nil {
		log.Fatal("failed to connect db:", err)
	}

	return db
}
