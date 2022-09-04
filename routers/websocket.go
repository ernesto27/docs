package routers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ernesto27/docs/interfaces"
	"github.com/ernesto27/docs/structs"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wss structs.WebsocketServer
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var broadcastDoc = make(chan structs.Doc)

func init() {
	wss = structs.WebsocketServer{
		Clients: make(map[int]structs.Client),
	}
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request, c *gin.Context, db interfaces.DocDB) {
	fmt.Println("Handle websocket connection")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgrade:", err)
		return
	}
	defer conn.Close()

	client := structs.Client{
		WebSocketConn: conn,
	}
	wss.Clients[10] = client

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error during message reading:", err)
			break
		}

		var command structs.Command
		err = json.Unmarshal(p, &command)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("COMMAND WS", command)

		switch command.Command {
		case "get-doc":
			fmt.Println("GET DOC")
			// TODO CHECK ID IS INT, PREVENT CRASH
			doc, err := db.GetDocByID(command.ID)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("DOC", doc)
			errWrite := wss.Clients[10].WebSocketConn.WriteJSON(doc)
			if errWrite != nil {
				log.Printf("error: %v", err)
				wss.Clients[10].WebSocketConn.Close()
				// delete(clients, wss.Clients[10].WebSocketConn)
			}

			break
		case "update-doc":
			fmt.Println("UPDATE DOC")
			break
		}
	}
}
