// Package chat handles chat logic
package chat

import (
	"chatter/server/internal/middleware"
	"chatter/server/internal/user"
	"encoding/json"
	"errors"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

const (
	maxConnections = 1000
	maxMessageSize = 512
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = 54 * time.Second
)

type messageType string

type Handler struct {
	service     *Service
	connections int32
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type sendMessageRequest struct {
	Message string `json:"message"`
}

type errorResponse struct {
	Message string `json:"message"`
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/chatroom", h.sendChatroomMessage)
	r.Get("/ws", h.readChatroomMessages)

	return r
}

func (h *Handler) sendChatroomMessage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	claimsInterface := r.Context().Value(middleware.UserKey)
	claims, ok := claimsInterface.(*user.CustomClaims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResponse{Message: "Unauthorized"})
		return
	}

	var req sendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Message: "Can't decode the JSON"})
		return
	}

	var m Message
	m.From = claims.UserID
	m.FromName = claims.Username
	m.Content = req.Message

	if err := h.service.SendChatroomMessage(r.Context(), &m); err != nil {
		if errors.Is(err, ErrNoMessage) || errors.Is(err, ErrMessageLimit) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse{Message: err.Error()})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse{Message: "Something went wrong"})
		return
	}

	w.WriteHeader(http.StatusCreated)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *Handler) readChatroomMessages(w http.ResponseWriter, r *http.Request) {
	if (atomic.LoadInt32(&h.connections)) >= maxConnections {
		http.Error(w, "Too many connections", http.StatusTooManyRequests)
		return
	}

	claimsInterface := r.Context().Value(middleware.UserKey)
	claims, ok := claimsInterface.(*user.CustomClaims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResponse{Message: "Unauthorized"})
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Unable to upgrade", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	atomic.AddInt32(&h.connections, 1)
	defer atomic.AddInt32(&h.connections, -1)

	conn.SetReadLimit(maxMessageSize)
	conn.SetWriteDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error { conn.SetWriteDeadline(time.Now().Add(pongWait)); return nil })

	h.service.Addclient(conn, claims.UserID, claims.Username)
	defer h.service.RemoveClient(conn, claims.UserID, claims.Username)

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			if _, _, err := conn.NextReader(); err != nil {
				break
			}
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
