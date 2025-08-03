// Package chat handles chat logic
package chat

import (
	"chatter/server/internal/middleware"
	"chatter/server/internal/user"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type Handler struct {
	service *Service
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
		json.NewEncoder(w).Encode(errorResponse{Message: "Bad request"})
		return
	}

	var m Message
	m.From = claims.UserID
	m.FromName = claims.Username
	m.Content = req.Message

	if err := h.service.SendChatroomMessage(r.Context(), &m); err != nil {
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
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Unable to upgrade", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	h.service.Addclient(conn)
	defer h.service.RemoveClient(conn)

	for {
		if _, _, err := conn.NextReader(); err != nil {
			break
		}
	}
}
