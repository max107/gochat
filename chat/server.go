package chat

import (
	"errors"
	"net/http"
)

type Server struct {
	Name     string              `json:"name"`
	Channels map[string]*Channel `json:"channels"`
	Users    []*User             `json:"users"`
	h        *Hub
}

func (self *Server) CreateChannel(name string) {
	self.Channels[name] = NewChannel(name)
}

func (self *Server) UpdateWsHandler(name string, w http.ResponseWriter, r *http.Request) {
	user, err := self.FindUser(name)
	if err != nil {
		user = NewUser(name)
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	c := NewConnection(ws, user)
	user.h.register <- c
	defer func() {
		user.h.unregister <- c
	}()
	go c.writer()
	c.reader()
}

func (self *Server) BroadcastToUsers(msg []byte) {
	go func() {
		for _, user := range self.Users {
			// Send message from system to user
		}
	}()
}

func (self *Server) BroadcastToChannels(msg []byte) {
	go func() {
		for name, channel := range self.Channels {
			// Send message from system to channels
		}
	}()
}

func (self *Server) FindUser(name string) (*User, error) {
	for _, user := range self.Users {
		if user.Name == name {
			return user, nil
		}
	}
	return &User{}, errors.New("User not found")
}

func NewServer() *Server {
	h := NewHub()
	go h.Run()
	return &Server{h: h}
}
