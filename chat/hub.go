package chat

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
