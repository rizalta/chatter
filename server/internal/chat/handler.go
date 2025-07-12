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
	}

	client.hub.register <- client

	go client.readPump()
	go client.writePump()
}
