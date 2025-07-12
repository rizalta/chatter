// Package chat handles the chat functions
package chat

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func HandleWS(hub *hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ws upgrade error: %v", err)
		return
	}

	client := &client{
		id:   uuid.New(),
		hub:  hub,
		conn: conn,
		send: make(chan []byte),
	}

	client.hub.register <- client

	go client.readPump()
	go client.writePump()
}
