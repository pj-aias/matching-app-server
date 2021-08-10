package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	mysql_password string
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func init() {
	f, err := os.Open("/run/secrets/mysql_password")
	if err != nil {
		log.Fatal("failed to open MySQL Password: ", err)
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

	mysql_password = string(buf)
}
