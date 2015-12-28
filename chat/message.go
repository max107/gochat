package chat

type Message struct {
	From string `json:"from"`
	Body []byte `json:"body"`
}
