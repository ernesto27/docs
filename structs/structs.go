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
	ID            int
	WebSocketConn *websocket.Conn
}

type WebsocketServer struct {
	Clients map[int]Client
}

type Command struct {
	Command string `json:"command"`
	Body    string `json:"body"`
	ID      int    `json:"id"`
}
