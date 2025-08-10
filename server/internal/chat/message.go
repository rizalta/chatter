package chat

import "time"

type Message struct {
	ID        string    `json:"id"`
	From      string    `json:"from"`
	FromName  string    `json:"fromName"`
	To        string    `json:"to,omitempty"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

type status string

const (
	statusJoined status = "joined"
	statusLeft   status = "left"
)

const (
	typePresence messageType = "presence"
	typeChat     messageType = "chat"
	typeUserList messageType = "user_list"
)

type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type PresenceMessage struct {
	Status status   `json:"status"`
	User   UserInfo `json:"user"`
}

type WSMessage struct {
	Type messageType `json:"type"`
	Data any         `json:"data"`
}
