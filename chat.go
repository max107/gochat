package main

import (
	"errors"
	"github.com/gorilla/websocket"
	"net/http"
)

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

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Message struct {
	From string `json:"from"`
	Body []byte `json:"body"`
}

type User struct {
	Name     string    `json:"name"`
	Messages []Message `json:"messages"`
	Channels []Channel `json:"channels"`
	h        *Hub
}

func (self *User) SendMessageToUser(to string, msg Message) bool {
	return true
}

func (self *User) SendMessageToChannel(to string, msg Message) bool {
	return true
}

type Channel struct {
	Name  string  `json:"name"`
	Users []*User `json:"users"`
	h     *Hub
}

type Server struct {
	Name     string    `json:"name"`
	Channels []Channel `json:"channels"`
	Users    []*User   `json:"users"`
	h        *Hub
}

func (self *Server) UpdateWsHandler(name string, w http.ResponseWriter, r *http.Request) {
	user, err := self.FindUser(name)
	if err != nil {
		user = self.NewUser(name)
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	c := &connection{
		send: make(chan []byte, 256),
		ws:   ws,
		user: user,
	}
	user.h.register <- c
	defer func() {
		user.h.unregister <- c
	}()
	go c.writer()
	c.reader()
}

func (self *Server) FindUser(name string) (*User, error) {
	for _, user := range self.Users {
		if user.Name == name {
			return user, nil
		}
	}
	return &User{}, errors.New("User not found")
}

func (self *Server) GetOrCreateUser(name string) *User {
	if user, err := self.FindUser(name); err == nil {
		return user
	}

	return self.NewUser(name)
}

func (self *Server) NewUser(name string) *User {
	h := NewHub()
	go h.Run()
	user := &User{Name: name, h: h}
	self.Users = append(self.Users, user)
	return user
}

func NewServer() *Server {
	h := NewHub()
	go h.Run()
	return &Server{h: h}
}

type Hub struct {
	// Registered connections.
	connections map[*connection]bool
	// Inbound messages from the connections.
	broadcast chan []byte
	// Register requests from the connections.
	register chan *connection
	// Unregister requests from connections.
	unregister chan *connection
}

func NewHub() *Hub {
	return &Hub{
		broadcast:   make(chan []byte),
		register:    make(chan *connection),
		unregister:  make(chan *connection),
		connections: make(map[*connection]bool),
	}
}

func (self *Hub) Run() {
	for {
		select {
		case c := <-self.register:
			self.connections[c] = true
		case c := <-self.unregister:
			if _, ok := self.connections[c]; ok {
				delete(self.connections, c)
				close(c.send)
			}
		case m := <-self.broadcast:
			for c := range self.connections {
				select {
				case c.send <- m:
				default:
					delete(self.connections, c)
					close(c.send)
				}
			}
		}
	}
}
