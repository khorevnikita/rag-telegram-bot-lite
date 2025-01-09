package main

import (
	"gopkg.in/telebot.v4"
	"gorag-telegram-bot/api"
	"gorag-telegram-bot/config"
	"gorag-telegram-bot/core"
	"gorag-telegram-bot/database"
	"log"
	"time"
)

func main() {
	config.LoadConfig("bot.yaml")
	database.Connect()
	defer database.CloseConnection()

	database.RedisConnect()
	defer database.CloseRedis()

	pref := telebot.Settings{
		Token:  config.AppConfig.BotToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}
	core.Bot = bot

	// Запуск HTTP-сервера
	router := api.NewRouter(bot)
	router.Listen()

	// Обработка сигнала завершения для корректного закрытия приложения
	/*c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c*/
}
