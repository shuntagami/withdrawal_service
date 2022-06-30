package main

import (
	"api/api"
	"api/db"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	dbConfig := db.Config{
		Charset:   os.Getenv("DB_CHARSET"),
		Collation: os.Getenv("DB_COLLATION"),
		Host:      os.Getenv("DB_HOST"),
		Name:      os.Getenv("DB_NAME"),
		Password:  os.Getenv("DB_PASSWORD"),
		Port:      os.Getenv("DB_PORT"),
		Username:  os.Getenv("DB_USERNAME"),
	}
	if err := db.Setup(&dbConfig); err != nil {
		log.Fatal(err)
	}
	r := gin.Default()
	r.GET("/_healthcheck", api.Healthcheck)
	r.POST("/histories", api.CreateHistory)
	if err := r.Run(":3000"); err != nil {
		log.Fatal(err)
	}
}
