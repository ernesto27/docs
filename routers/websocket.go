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

const (
	INIT             = "init"
	GET_DOC          = "get-doc"
	UPDATE_DOC_BODY  = "update-doc-body"
	UPDATE_DOC_TITLE = "update-doc-title"
	USERS_CONNECTED  = "users-connected"
	ERROR            = "error"
)

var wss structs.WebsocketServer
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var chanDocID = make(chan int)
var chanUserID = make(chan int)

type DocBase struct {
	doc           structs.Doc
	currentID     int
	typeBroadcast string
}

var broadcastDoc = make(chan DocBase)

func init() {
	wss = structs.WebsocketServer{
		Clients: make(map[int]structs.Client),
	}
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request, c *gin.Context, db interfaces.DocDB) {
	fmt.Println("Handle websocket connection")
	rand.Seed(time.Now().UnixNano())
	userID := rand.Int()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgrade:", err)
		return
	}
	defer conn.Close()

	client := structs.Client{
		ID:            userID,
		WebSocketConn: conn,
	}
	wss.Clients[userID] = client

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			docID := wss.Clients[userID].DocID
			delete(wss.Clients, userID)
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
		case GET_DOC:
			fmt.Println("COMMAND WS", command)
			NewUserConnected(command, userID)
			// TODO CHECK ID IS INT, PREVENT CRASH
			getDoc(db, command, userID, INIT)
			break
		case UPDATE_DOC_BODY:
			updateDocBody(db, command, userID)
			break
		case UPDATE_DOC_TITLE:
			updateDocTitle(db, command, userID)
			break
		}
	}
}

func updateDocTitle(db interfaces.DocDB, command structs.Command, userID int) {
	fmt.Println("UPDATE DOC TITLE")
	rows, err := db.UpdateDocTitleByID(command.ID, command.Title)
	if rows == 0 || err != nil {
		fmt.Println(err)
	} else {
		doc, errDoc := db.GetDocByID(command.ID)
		if errDoc != nil {
			fmt.Println(errDoc)
		} else {
			broadcastDoc <- DocBase{
				doc:       doc,
				currentID: userID,
			}
		}
	}
}

func updateDocBody(db interfaces.DocDB, command structs.Command, userID int) {
	fmt.Println("UPDATE DOC")
	rows, err := db.UpdateDocBodyByID(command.ID, command.Body)
	if rows == 0 || err != nil {
		fmt.Println(err)
		chanUserID <- userID
	} else {
		doc, err := db.GetDocByID(command.ID)
		if err != nil {
			fmt.Println(err)
		} else {
			broadcastDoc <- DocBase{
				doc:       doc,
				currentID: userID,
			}
		}
	}
}

func getDoc(db interfaces.DocDB, command structs.Command, userID int, typeBroadcast string) {
	fmt.Println("GET DOC")
	doc, err := db.GetDocByID(command.ID)
	if err != nil {
		fmt.Println(err)
		chanUserID <- userID
	} else {
		broadcastDoc <- DocBase{
			doc:           doc,
			currentID:     userID,
			typeBroadcast: typeBroadcast,
		}
	}
}

func NewUserConnected(command structs.Command, userID int) {
	currentClient := wss.Clients[userID]
	currentClient.DocID = command.ID
	wss.Clients[userID] = currentClient
	chanDocID <- command.ID
}

func BroadcastDocByID() {
	for {
		docBase := <-broadcastDoc
		fmt.Println("BROADCAST DOC", docBase)
		responseDocByID := structs.ResponseDocByID{
			Command: GET_DOC,
			Doc:     docBase.doc,
		}

		if docBase.typeBroadcast == INIT {
			client := wss.Clients[docBase.currentID]
			writeJSON(client, responseDocByID)

		} else {
			for key, client := range wss.Clients {
				if key != docBase.currentID {
					writeJSON(client, responseDocByID)
				}
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
			Command: USERS_CONNECTED,
			Users:   users,
		}

		for _, client := range wss.Clients {
			if client.DocID == docID {
				writeJSON(client, responseUsersConnected)
			}
		}
	}
}

func BroadcastError() {
	for {
		userID := <-chanUserID
		fmt.Println("BROADCAST ERROR")
		responseError := structs.ResponseError{
			Command: ERROR,
			Error:   "Error ocurred, please try again",
		}

		client := wss.Clients[userID]
		writeJSON(client, responseError)
	}
}

func writeJSON(client structs.Client, response interface{}) {
	errWebsocket := client.WebSocketConn.WriteJSON(response)
	if errWebsocket != nil {
		log.Printf("error: %v", errWebsocket)
		client.WebSocketConn.Close()
		delete(wss.Clients, client.ID)
	}
}
