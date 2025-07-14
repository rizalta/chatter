package chat

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writWait       = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4 * 1024 // 4 KB
)

var (
	space   = " "
	newLine = "\n"
)

type client struct {
	id   string
	conn *websocket.Conn
	hub  *hub
	send chan *Message
}

func (c *client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(int64(maxMessageSize))
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(appData string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("readPump error: %v", err)
			}
			break
		}

		var incoming IncomingMessage
		if err := json.Unmarshal(data, &incoming); err != nil {
			log.Printf("readPump unmarshall err: %v", err)
			continue
		}

		message := &Message{
			Content:   strings.TrimSpace(strings.ReplaceAll(incoming.Content, newLine, space)),
			To:        "chatroom",
			From:      c.id,
			Timestamp: time.Now().Unix(),
		}

		c.hub.broadcast <- message
	}
}

func (c *client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(pingPeriod))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Printf("writePump marshall error: %v", err)
				continue
			}

			w.Write(data)
			n := len(c.send)
			for range n {
				m := <-c.send
				w.Write([]byte(newLine))
				data, err := json.Marshal(m)
				if err != nil {
					log.Printf("writePump marshall error: %v", err)
					continue
				}
				w.Write(data)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(pingPeriod))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
