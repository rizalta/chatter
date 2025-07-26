package user

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

type request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

type errorResponse struct {
	Message string `json:"message"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/login", h.handleLogin)
	r.Post("/register", h.handleRegister)

	return r
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserNotFound), errors.Is(err, ErrInvalidCredentials):
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(errorResponse{Message: err.Error()})
		default:
			log.Printf("internal server error during login, %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loginResponse{Token: token})
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid reuest", http.StatusBadRequest)
		return
	}

	err := h.service.Register(r.Context(), req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, ErrUsernameLength),
			errors.Is(err, ErrUsernameStart),
			errors.Is(err, ErrUsernameContains),
			errors.Is(err, ErrPasswordLength),
			errors.Is(err, ErrPasswordDigit),
			errors.Is(err, ErrPasswordLowercase),
			errors.Is(err, ErrPasswordUppercase),
			errors.Is(err, ErrPasswordSpecial):
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse{Message: err.Error()})

		case errors.Is(err, ErrUsernameAlreadyExists):
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(errorResponse{Message: err.Error()})

		default:
			log.Printf("internal server error during registration, %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
