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
