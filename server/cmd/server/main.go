package main

import (
	"chatter/server/config"
	"chatter/server/internal/chat"
	"chatter/server/internal/database"
	"chatter/server/internal/middleware"
	"chatter/server/internal/user"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	config, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config, %v", err)
	}

	ctx := context.Background()

	db, err := database.NewClient(ctx, config.RedisAddr)
	if err != nil {
		log.Fatalf("Error connecting to db, %v", err)
	}

	router := chi.NewRouter()

	router.Use(chimiddleware.Logger)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	}))

	userRepo := database.NewUserRepo(db)
	userService := user.NewService(userRepo, config.JWTPrivateKey)
	userHandler := user.NewHandler(userService)

	router.Mount("/api/user", userHandler.Routes())

	chatRepo := database.NewChatRepo(db)
	chatService := chat.NewService(chatRepo)
	chatHandler := chat.NewHandler(chatService)

	router.Group(func(r chi.Router) {
		r.Use(middleware.Auth(config.JWTPublicKey))
		r.Mount("/api/chat", chatHandler.Routes())
	})

	go chatService.Listen(ctx)

	log.Printf("Running server on port: %s", config.ServerPort)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.ServerPort), router))
}
