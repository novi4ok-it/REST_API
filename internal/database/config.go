package database

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
)

type Config struct {
	Port string
	DB   DBConfig
}

type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
	SSLMode  string
}

var AppConfig Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Не удалось загрузить .env файл, используются системные переменные окружения")
	}

	viper.SetConfigFile("config.yml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Не удалось загрузить config.yml: %v", err)
	}

	AppConfig = Config{
		Port: viper.GetString("PORT"),
		DB: DBConfig{
			Username: viper.GetString("db.username"),
			Password: os.Getenv("DB_PASSWORD"),
			Host:     viper.GetString("db.host"),
			Port:     viper.GetString("db.port"),
			DBName:   viper.GetString("db.dbname"),
			SSLMode:  viper.GetString("db.sslmode"),
		},
	}
}
