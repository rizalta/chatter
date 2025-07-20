package main

import (
	"chatter/server/config"
	"chatter/server/internal/chat"
	"chatter/server/internal/database"
	"chatter/server/internal/middleware"
	"chatter/server/internal/user"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)

	config, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config, %v", err)
	}

	db, err := database.NewClient(config.RedisAddr)
	if err != nil {
		log.Fatalf("Error connecting to db, %v", err)
	}

	userRepo := database.NewUserRepo(db)
	userService := user.NewService(userRepo, config.JWTPrivateKey)
	userHandler := user.NewHandler(userService)

	r.Mount("/api/user/", userHandler.Routes())

	hub := chat.NewHub()
	go hub.Run()

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(config.JWTPublicKey))

		r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			chat.HandleWS(hub, w, r)
		})
	})

	log.Printf("Running server on port: %s", config.ServerPort)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.ServerPort), r))
}
