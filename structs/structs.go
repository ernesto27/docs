package structs

import (
	"time"

	"github.com/gorilla/websocket"
)

type Doc struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ResponseApi struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Client struct {
	ID            int             `json:"id"`
	WebSocketConn *websocket.Conn `json:"-"`
	DocID         int
}

type WebsocketServer struct {
	Clients map[int]Client
}

type Command struct {
	Command string `json:"command"`
	Body    string `json:"body"`
	Title   string `json:"title"`
	ID      int    `json:"id"`
}

type ResponseDocByID struct {
	Command string `json:"command"`
	Doc     Doc    `json:"doc"`
}

type ResponseUsersConnected struct {
	Command string   `json:"command"`
	Users   []Client `json:"users"`
}

type ResponseError struct {
	Command string `json:"command"`
	Error   string `json:"error"`
}
