package chat

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

func NewUser(name string) *User {
	h := NewHub()
	go h.Run()
	return &User{Name: name, h: h}
}
