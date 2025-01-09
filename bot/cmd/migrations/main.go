package main

import (
	"flag"
	"fmt"
	"github.com/pressly/goose/v3"
	"gorag-telegram-bot/config"
	"gorag-telegram-bot/database"
	"log"
)

func main() {
	// Загружаем конфигурацию и подключаемся к базе данных
	config.LoadConfig("bot.yaml")
	database.Connect()
	gormDB := database.DB

	// Получаем *sql.DB из *gorm.DB для совместимости с goose
	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatalf("Failed to get *sql.DB from *gorm.DB: %v", err)
	}
	// Определяем путь к миграциям
	migrationsDir := "./database/migrations"

	// Обработка аргументов командной строки
	command := flag.String("command", "up", "Specify 'up' for applying migrations or 'down' for rolling back")
	flag.Parse()

	// Выполнение команды
	switch *command {
	case "up":
		if err := goose.Up(sqlDB, migrationsDir); err != nil {
			log.Fatalf("Failed to apply migrations: %v", err)
		}
		fmt.Println("Migrations applied successfully.")
	case "down":
		if err := goose.Down(sqlDB, migrationsDir); err != nil {
			log.Fatalf("Failed to roll back migrations: %v", err)
		}
		fmt.Println("Migrations rolled back successfully.")
	default:
		fmt.Println("Invalid command. Use 'up' or 'down'.")
	}
}
