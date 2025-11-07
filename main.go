package main

import (
	"chat-server/internal/config"
	"chat-server/internal/db"
	handlers "chat-server/internal/handler"
	"log"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()
	database := db.Connect(cfg)

	http.Handle("/webhook", handlers.WebhookHandler(database))
	url := ":8080"
	log.Printf("Cервер запущен на %s", url)
	log.Fatal(http.ListenAndServe(url, nil))
}
