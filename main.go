package main

import (
	"log"
	"net/http"

	"github.com/ernesto27/docs/db"
	"github.com/ernesto27/docs/routers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	myDb := db.Mysql{}
	myDb.New()

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.POST("/docs/create", func(c *gin.Context) {
		routers.CreateDoc(&myDb, c)
	})

	r.Run()
}
