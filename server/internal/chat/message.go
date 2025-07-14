package chat

type Message struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

type IncomingMessage struct {
	Content string `json:"content"`
	To      string `json:"to"`
}
