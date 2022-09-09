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
var chanDocID = make(chan int)
var chanUserID = make(chan int)

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
		ID:            id,
		WebSocketConn: conn,
	}
	wss.Clients[id] = client

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			docID := wss.Clients[id].DocID
			delete(wss.Clients, id)
			chanDocID <- docID
			log.Println("Error during message reading:", err)
			break
		}

		var command structs.Command
		err = json.Unmarshal(p, &command)
		if err != nil {
			log.Fatal(err)
		}

		switch command.Command {
		case "get-doc":
			fmt.Println("COMMAND WS", command)
			NewUserConnected(id, command)
			// TODO CHECK ID IS INT, PREVENT CRASH
			getDoc(db, id, command)
			break
		case "update-doc-body":
			updateDocBody(db, id, command)
			break
		case "update-doc-title":
			updateDocTitle(db, command)
			break
		}
	}
}

func updateDocTitle(db interfaces.DocDB, command structs.Command) {
	fmt.Println("UPDATE DOC TITLE")
	rows, err := db.UpdateDocTitleByID(command.ID, command.Title)
	if rows == 0 || err != nil {
		fmt.Println(err)
	} else {
		doc, errDoc := db.GetDocByID(command.ID)
		if errDoc != nil {
			fmt.Println(errDoc)
		} else {
			broadcastDoc <- doc
		}
	}
}

func updateDocBody(db interfaces.DocDB, id int, command structs.Command) {
	fmt.Println("UPDATE DOC")
	rows, err := db.UpdateDocBodyByID(command.ID, command.Body)
	if rows == 0 || err != nil {
		fmt.Println(err)
		chanUserID <- id
	} else {
		doc, err := db.GetDocByID(command.ID)
		if err != nil {
			fmt.Println(err)
		} else {
			broadcastDoc <- doc
		}
	}
}

func getDoc(db interfaces.DocDB, id int, command structs.Command) {
	fmt.Println("GET DOC")
	doc, err := db.GetDocByID(command.ID)
	if err != nil {
		fmt.Println(err)
		chanUserID <- id
	} else {
		broadcastDoc <- doc
	}
}

func NewUserConnected(id int, command structs.Command) {
	currentClient := wss.Clients[id]
	currentClient.DocID = command.ID
	wss.Clients[id] = currentClient
	chanDocID <- command.ID
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
				delete(wss.Clients, client.ID)
			}
		}
	}
}

func BroadcastUsersConnected() {
	for {
		docID := <-chanDocID
		fmt.Println("BROADCAST USERS CONNECTED")

		users := []structs.Client{}

		for _, client := range wss.Clients {
			if client.DocID == docID {
				users = append(users, client)
			}
		}
		responseUsersConnected := structs.ResponseUsersConnected{
			Command: "users-connected",
			Users:   users,
		}

		for _, client := range wss.Clients {
			if client.DocID == docID {
				err := client.WebSocketConn.WriteJSON(responseUsersConnected)
				if err != nil {
					log.Printf("error: %v", err)
					client.WebSocketConn.Close()
					delete(wss.Clients, client.ID)
				}
			}
		}
	}
}

func BroadcastError() {
	for {
		userID := <-chanUserID
		fmt.Println("BROADCAST ERROR")
		responseError := structs.ResponseError{
			Command: "error",
			Error:   "Error ocurred, please try again",
		}

		client := wss.Clients[userID]
		errWebsocket := client.WebSocketConn.WriteJSON(responseError)
		if errWebsocket != nil {
			log.Printf("error: %v", errWebsocket)
			client.WebSocketConn.Close()
			delete(wss.Clients, client.ID)
		}
	}
}
