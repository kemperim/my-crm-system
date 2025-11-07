package models

import (
	"time"
)

type Message struct {
	ID        int    `gorm:"primaryKey"`
	ChatID    string `gorm:"index"`
	Sender    string
	Text      string
	Timestamp time.Time
}
