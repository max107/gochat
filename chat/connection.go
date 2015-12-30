package chat

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type connection struct {
	// The websocket connection.
	ws *websocket.Conn
	// Buffered channel of outbound messages.
	send chan []byte
	user *User
}

func (c *connection) reader() {
	for {
		_, body, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		c.user.h.broadcast <- body
	}
	c.ws.Close()
}

func (c *connection) writer() {
	for message := range c.send {
		err := c.ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
	c.ws.Close()
}

func NewConnection(ws *websocket.Conn, user *User) *connection {
	return &connection{send: make(chan []byte, 256), ws: ws, user: user}
}
