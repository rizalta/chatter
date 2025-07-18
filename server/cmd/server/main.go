package main

import (
	"chatter/server/config"
	"chatter/server/internal/chat"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	config, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config, %v", err)
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})

	hub := chat.NewHub()
	go hub.Run()

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		chat.HandleWS(hub, w, r)
	})

	log.Printf("Running server on port: %s", config.ServerPort)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.ServerPort), r))
}
