package chat

type Channel struct {
	Name  string  `json:"name"`
	Users []*User `json:"users"`
	h     *Hub
}

func NewChannel(name string) *Channel {
	h := NewHub()
	go h.Run()
	return &Channel{Name: name, h: h}
}
