package db

import (
	"chat-server/internal/config"
	"chat-server/internal/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) *gorm.DB {

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.Message{})
	log.Println("✅ Подключение к PostgreSQL успешно")
	return db

}

var DB *gorm.DB

func InitDB() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if host == "" || user == "" || dbname == "" {
		log.Fatal(" Не хватает параметров БД в .env")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf(" Ошибка подключения к базе: %v", err)
	}

	if err := db.AutoMigrate(&models.Message{}); err != nil {
		log.Fatalf(" Ошибка миграции: %v", err)
	}

	DB = db
	fmt.Println(" Подключено к PostgreSQL и выполнена миграция")
}
