// Package chat handles chat logic
package chat

import (
	"chatter/server/internal/middleware"
	"chatter/server/internal/user"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
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
	return r
}

func (h *Handler) sendChatroomMessage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	claimsInterface := r.Context().Value(middleware.UserKey)
	fmt.Printf("claimsInterface: %v\n", claimsInterface)
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
	m.Content = req.Message

	if err := h.service.SendChatroomMessage(r.Context(), &m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse{Message: "Something went wrong"})
		return
	}

	w.WriteHeader(http.StatusOK)
}
