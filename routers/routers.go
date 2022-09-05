package routers

import (
	"fmt"
	"net/http"

	"github.com/ernesto27/docs/interfaces"
	"github.com/ernesto27/docs/structs"
	"github.com/gin-gonic/gin"
)

func CreateDoc(db interfaces.DocDB, c *gin.Context) {
	doc := structs.Doc{
		Title: c.PostForm("title"),
		Body:  c.PostForm("body"),
	}
	id, err := db.CreateDoc(doc)

	if err != nil {
		c.JSON(http.StatusOK, structs.ResponseApi{
			Status:  "error",
			Message: "error creating doc",
		})
		return
	}
	c.Redirect(http.StatusFound, "/docs/"+fmt.Sprint(id))
}
