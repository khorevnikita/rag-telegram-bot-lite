package handlers

import (
	tb "gopkg.in/telebot.v4"
	"gorag-telegram-bot/config"
	"gorag-telegram-bot/core"
	"gorag-telegram-bot/services"
)

type ReportHandlers struct {
	Bot *tb.Bot
}

func GetNewReportHandlers() ReportHandlers {
	return ReportHandlers{
		Bot: core.Bot,
	}
}

func (h ReportHandlers) RegisterCommands() {
	h.Bot.Handle(&tb.InlineButton{Unique: "report_users"}, h.usersList)
	h.Bot.Handle(&tb.InlineButton{Unique: "report_subscriptions"}, h.subscriptionsList)

	if config.AppConfig.Modules.Likes.Enabled {
		h.Bot.Handle(&tb.InlineButton{Unique: "report_dislike_messages"}, h.dislikeMessages)
	}
}

func (h ReportHandlers) usersList(c tb.Context) error {
	_ = c.Respond()
	users, err := services.GetUsers()
	if err != nil {
		return c.Send(err.Error())
	}
	questions, err := services.GetQuestions()
	if err != nil {
		return c.Send(err.Error())
	}
	answers, err := services.GetAnswers()
	if err != nil {
		return c.Send(err.Error())
	}

	// Экспорт в Excel
	buf, err := services.ExportUsersToExcel(users, questions, answers)
	if err != nil {
		return c.Send("Failed to generate report: " + err.Error())
	}

	// Отправляем файл пользователю
	file := &tb.Document{
		File:     tb.File{FileReader: buf},
		FileName: "users_report.xlsx",
	}
	return c.Send(file)
}

func (h ReportHandlers) subscriptionsList(c tb.Context) error {
	_ = c.Respond()
	subscriptions, err := services.GetSubscriptionReport()
	if err != nil {
		return c.Send(err.Error())
	}

	// Экспорт в Excel
	buf, err := services.ExportSubscriptions(subscriptions)
	if err != nil {
		return c.Send("Failed to generate report: " + err.Error())
	}

	// Отправляем файл пользователю
	file := &tb.Document{
		File:     tb.File{FileReader: buf},
		FileName: "subscription_report.xlsx",
	}
	return c.Send(file)
}

func (h ReportHandlers) dislikeMessages(c tb.Context) error {
	_ = c.Respond()
	messages, err := services.GetDislikedMessages()
	if err != nil {
		return c.Send(err.Error())
	}

	// Экспорт в Excel
	buf, err := services.ExportDislikedMessages(messages)
	if err != nil {
		return c.Send("Failed to generate report: " + err.Error())
	}

	// Отправляем файл пользователю
	file := &tb.Document{
		File:     tb.File{FileReader: buf},
		FileName: "dislikes_report.xlsx",
	}
	return c.Send(file)
}
