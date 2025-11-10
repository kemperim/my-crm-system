package handlers

import (
	"chat-server/internal/models"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"
)

func WebhookHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var data map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		typeWebhook, _ := data["typeWebhook"].(string)

		log.Println(typeWebhook)

		if typeWebhook == "incomingMessageReceived" {

			senderData := map[string]interface{}{}
			if sd, ok := data["senderData"].(map[string]interface{}); ok {
				senderData = sd
			} else {
				log.Println("senderData не найден")
			}

			text := ""
			if messageData, ok := data["messageData"].(map[string]interface{}); ok {
				if ext, ok := messageData["extendedTextMessageData"].(map[string]interface{}); ok {
					if t, ok := ext["text"].(string); ok {
						text = t
					}
				}
			}

			message := models.Message{
				ChatID:    senderData["chatId"].(string),
				Sender:    senderData["chatName"].(string),
				Text:      text,
				Timestamp: time.Now(),
			}
			log.Println(message)

			if err := db.Create(&message).Error; err != nil {
				http.Error(w, "Failed to save message", http.StatusInternalServerError)
				log.Println(" Ошибка сохранения в БД:", err)
				return
			}

			log.Printf(" Сообщение сохранено: ID=%d, ChatID=%s, Sender=%s, Text=%s\n",
				message.ID, message.ChatID, message.Sender, message.Text,
			)

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))

		} else if typeWebhook == "incomingCall" && data["status"] != "offer" {
			log.Println(data)

			call := models.Call{
				ChatID:    data["from"].(string),
				Sender:    data["from"].(string),
				Status:    data["status"].(string),
				Timestamp: time.Now(),
			}
			if err := db.Create(&call).Error; err != nil {
				http.Error(w, "Failed to save message", http.StatusInternalServerError)
				log.Println(" Ошибка сохранения в БД:", err)
				return
			}
		}
	}

}
