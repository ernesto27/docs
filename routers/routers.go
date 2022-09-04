package routers

import (
	"net/http"

	"github.com/ernesto27/docs/interfaces"
	"github.com/ernesto27/docs/structs"
	"github.com/gin-gonic/gin"
)

func CreateDoc(db interfaces.DocDB, c *gin.Context) {
	// Validate form params
	title := c.PostForm("title")
	body := c.PostForm("body")

	if title == "" || body == "" {
		c.JSON(http.StatusOK,
			structs.ResponseApi{
				Status:  "error",
				Message: "title or body is empty",
			},
		)
		return
	}

	doc := structs.Doc{
		Title: c.PostForm("title"),
		Body:  c.PostForm("body"),
	}
	err := db.CreateDoc(doc)

	if err != nil {
		c.JSON(http.StatusOK, structs.ResponseApi{
			Status:  "error",
			Message: "error creating doc",
		})
		return
	}

	c.JSON(http.StatusOK, structs.ResponseApi{
		Status:  "success",
		Message: "success created doc",
	})
}
