package config

import (
	"RestAPI/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDB() {
	dsn := "host=localhost user=postgres password=mysecretpassword dbname=postgres port=5438 sslmode=disable"
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
