package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"chat-server/internal/models"

	"gorm.io/gorm"
)

// WebhookHandler возвращает http.HandlerFunc с доступом к БД
func WebhookHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var data WebhookData
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			log.Println("❌ Ошибка парсинга webhook:", err)
			return
		}

		// Логируем входящие данные для проверки
		log.Printf("📩 Получено сообщение: ChatID=%s, Sender=%s, Text=%s\n",
			data.SenderData.ChatID,
			data.SenderData.ChatName,
			data.MessageData.ExtendedTextMessageData.Text,
		)

		// Создаем объект модели
		message := models.Message{
			ChatID:    data.SenderData.ChatID, // теперь точно экспортируемое поле
			Sender:    data.SenderData.ChatName,
			Text:      data.MessageData.ExtendedTextMessageData.Text,
			Timestamp: time.Now(),
		}

		// Сохраняем в базу
		if err := db.Create(&message).Error; err != nil {
			http.Error(w, "Failed to save message", http.StatusInternalServerError)
			log.Println("❌ Ошибка сохранения в БД:", err)
			return
		}

		// Логируем сохранённое сообщение с ID
		log.Printf("💾 Сообщение сохранено: ID=%d, ChatID=%s, Sender=%s, Text=%s\n",
			message.ID, message.ChatID, message.Sender, message.Text,
		)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))

		// Для дебага в консоль
		fmt.Printf("%+v\n", message)
	}
}

// WebhookData структура входящего webhook
type WebhookData struct {
	TypeWebhook  string `json:"typeWebhook"`
	InstanceData struct {
		IDInstance   int64  `json:"idInstance"`
		Wid          string `json:"wid"`
		TypeInstance string `json:"typeInstance"`
	} `json:"instanceData"`
	Timestamp  int64  `json:"timestamp"`
	IDMessage  string `json:"idMessage"`
	SenderData struct {
		ChatID            string `json:"chatId"`
		ChatName          string `json:"chatName"`
		Sender            string `json:"sender"`
		SenderName        string `json:"senderName"`
		SenderContactName string `json:"senderContactName"`
	} `json:"senderData"`
	MessageData struct {
		TypeMessage             string `json:"typeMessage"`
		ExtendedTextMessageData struct {
			Text string `json:"text"`
		} `json:"extendedTextMessageData"`
	} `json:"messageData"`
}
