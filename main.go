package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
	myDb.New(os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_NAME"))

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.POST("/docs/create", func(c *gin.Context) {
		routers.CreateDoc(&myDb, c)
	})

	r.GET("/docs/:id", func(c *gin.Context) {
		id := c.Param("id")
		fmt.Println(id)
		c.HTML(http.StatusOK, "doc.html", nil)
	})

	r.GET("/ws", func(c *gin.Context) {
		routers.WebsocketHandler(c.Writer, c.Request, c, &myDb)
	})

	r.Run()
}
