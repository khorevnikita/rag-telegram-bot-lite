package main

import (
	"fmt"
	"gorag-telegram-bot/cmd/seeder/seeds"
	"gorag-telegram-bot/config"
	"gorag-telegram-bot/database"
	"os"
)

func main() {
	// Загружаем конфигурацию и подключаемся к базе данных
	config.LoadConfig("bot.yaml")
	database.Connect()
	db := database.DB // Используется db, чтобы избежать ошибок

	// Получаем аргументы
	args := os.Args[1:] // Пропускаем имя программы (os.Args[0])

	// Проверяем наличие аргументов
	if len(args) == 0 {
		fmt.Println("No arguments provided")
		return
	}

	// Обрабатываем аргументы
	switch args[0] {
	case "questions":
		seeds.SeedForm(db)
	case "contexts":
		seeds.SeedContexts(db)
	case "payment_providers":
		seeds.SeedPaymentProvider(db)
	default:
		fmt.Println("Unknown argument:", args[0])
	}
}
