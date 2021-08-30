package db

import (
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	gormMySQL "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	database *gorm.DB
)

func TestInsert(user User) {
	database.Create(&user)
}

func ConnectDB(dbName string) (*gorm.DB, error) {
	tz, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic("Failed to Parse TimeZone")
	}

	config := mysql.Config{
		User:      "root",
		Passwd:    getPassword(),
		Net:       "tcp",
		Addr:      "db" + ":3306",
		DBName:    dbName,
		ParseTime: true,
		Loc:       tz,
	}
	dsn := config.FormatDSN()

	conn := gormMySQL.Open(dsn)
	return gorm.Open(conn, &gorm.Config{})
}

func init() {
	dbName := "service"

	env := os.Getenv("SERVICE_ENV")

	if env == "test" {
		dbName = "test"
	}

	var err error
	database, err = ConnectDB(dbName)

	if err != nil {
		panic("Failed to open MySQL Connection")
	}

	database.Logger.LogMode(logger.Info)
	autoMigrate(database)
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
