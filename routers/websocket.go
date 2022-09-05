package routers

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

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

	rand.Seed(time.Now().UnixNano())
	id := rand.Int()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgrade:", err)
		return
	}
	defer conn.Close()

	client := structs.Client{
		WebSocketConn: conn,
	}
	wss.Clients[id] = client

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
				// TODO SEND ERROR MESSAGE
				fmt.Println(err)
			}
			// TODO SEND ONLY TO USER CONNECT FIRST TIME
			broadcastDoc <- doc

			break
		case "update-doc":
			fmt.Println("UPDATE DOC")
			rows, err := db.UpdateDocByID(command.ID, command.Body)
			if rows == 0 || err != nil {
				// TODO SEND ERROR MESSAGE
				fmt.Println(err)
			}

			doc, err := db.GetDocByID(command.ID)
			if err != nil {
				// TODO SEND ERROR MESSAGE
				fmt.Println(err)
			}
			broadcastDoc <- doc

			break
		}
	}
}

func BroadcastDocByID() {
	for {
		doc := <-broadcastDoc
		fmt.Println("BROADCAST DOC", doc)
		responseDocByID := structs.ResponseDocByID{
			Command: "get-doc",
			Doc:     doc,
		}
		for _, client := range wss.Clients {
			err := client.WebSocketConn.WriteJSON(responseDocByID)
			if err != nil {
				log.Printf("error: %v", err)
				client.WebSocketConn.Close()
				// delete(clients, wss.Clients[1].WebSocketConn)
			}
		}
	}
}
