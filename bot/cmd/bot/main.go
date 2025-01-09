package main

import (
	"gopkg.in/telebot.v4"
	"gorag-telegram-bot/config"
	"gorag-telegram-bot/core"
	"gorag-telegram-bot/database"
	"gorag-telegram-bot/handlers"
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

	_ = GetNewRAGBot()
	core.Bot.Start()
	log.Println("Бот запущен.")
}

type RAGBot struct {
	CommandHandler handlers.CommandHandlers
	LikesHandler   handlers.LikeHandlers
	InputHandler   handlers.InputHandlers
	FormHandler    handlers.FormHandlers
	ReportHandler  handlers.ReportHandlers
	CommandPrompt  handlers.CommandPromptHandlers
	BillingHandler handlers.BillingHandlers
	MenuHandler    handlers.MenuHandlers
}

func GetNewRAGBot() RAGBot {
	formHandler := handlers.GetNewFormHandlers()

	s := RAGBot{
		CommandHandler: handlers.GetNewCommandHandler(),
		LikesHandler:   handlers.GetNewLikeHandlers(),
		InputHandler:   handlers.GetNewInputHandlers(formHandler),
		FormHandler:    formHandler,
		ReportHandler:  handlers.GetNewReportHandlers(),
		CommandPrompt:  handlers.GetNewCommandPromptHandlers(),
		BillingHandler: handlers.GetBillingHandlers(),
		MenuHandler:    handlers.GetMenuHandlers(),
	}

	s.registerHandlers()
	s.bootEvents()
	s.setMainCommands()

	return s
}

func (s RAGBot) registerHandlers() {
	s.CommandHandler.RegisterCommands()
	s.LikesHandler.RegisterCommands()
	s.InputHandler.RegisterCommands()
	s.FormHandler.RegisterCommands()
	s.ReportHandler.RegisterCommands()
	s.CommandPrompt.RegisterCommands()
	s.BillingHandler.RegisterCommands()
	s.MenuHandler.RegisterCommands()
}

func (s RAGBot) bootEvents() {
	// some events to boot
}

func (s RAGBot) setMainCommands() {
	menuCommands := []telebot.Command{}
	for _, cmd := range config.AppConfig.Commands {
		if cmd.Enabled && cmd.ShowInMenu {
			menuCommands = append(menuCommands, telebot.Command{
				Text:        cmd.Name,
				Description: cmd.Label,
			})
		}
	}
	if err := core.Bot.SetCommands(menuCommands); err != nil {
		log.Fatalf("Error setting commands: %v", err)
	}
}
