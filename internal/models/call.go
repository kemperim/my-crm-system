package models

import "time"

type Call struct {
	ID        int    `gorm:"primaryKey"`
	ChatID    string `gorm:"index"`
	Sender    string
	Status    string
	Timestamp time.Time
}
