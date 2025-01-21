package config

import (
	"RestAPI/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDB() {
	dbConfig := AppConfig.DB
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		dbConfig.Host, dbConfig.Username, dbConfig.Password, dbConfig.DBName, dbConfig.Port, dbConfig.SSLMode,
	)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}

	// Автоматическая миграция
	err = DB.AutoMigrate(
		&models.User{},
		&models.TodoList{},
		&models.Task{},
	)
	if err != nil {
		log.Fatalf("Не удалось выполнить миграцию: %v", err)
	}
}
