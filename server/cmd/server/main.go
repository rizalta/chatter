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

	config := config.Load()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})

	hub := chat.NewHub()
	go hub.Run()

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		chat.HandleWS(hub, w, r)
	})

	log.Printf("Running server on port: %s", config.Port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.Port), r))
}
