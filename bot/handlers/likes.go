package handlers

import (
	"fmt"
	tb "gopkg.in/telebot.v4"
	"gorag-telegram-bot/config"
	"gorag-telegram-bot/core"
	"gorag-telegram-bot/services"
)

type LikeHandlers struct {
	likeConfig config.LikesModule
}

func GetNewLikeHandlers() LikeHandlers {
	likesConfig := config.AppConfig.Modules.Likes
	return LikeHandlers{
		likeConfig: likesConfig,
	}
}

func (h LikeHandlers) RegisterCommands() {
	if h.likeConfig.Enabled {
		core.Bot.Handle(&tb.InlineButton{Unique: "like_button"}, h.likeHandler)
		core.Bot.Handle(&tb.InlineButton{Unique: "dislike_button"}, h.dislikeHandler)
	}
}

func (h LikeHandlers) likeHandler(c tb.Context) error {
	_ = c.Respond()
	// Находим сообщение по TelegramMessageID
	messageID := c.Message().ID

	message, err := services.FindByTelegramID(messageID)
	if err != nil {
		return c.Respond(&tb.CallbackResponse{Text: "Сообщение не найдено."})
	}

	err = services.SetLiked(message)
	if err != nil {
		return c.Respond(&tb.CallbackResponse{Text: "Ошибка при обновлении оценки."})
	}

	// Отправляем уведомление в сервис ИИ
	aiClient := services.AIClient{}
	if err := aiClient.LikeQuestion(*message.AIMessageID); err != nil {
		fmt.Printf("Ошибка при отправке оценки в ИИ сервис: %v\n", err)
	}

	services.ClearSourceMessage(c)

	return c.Respond(&tb.CallbackResponse{Text: h.likeConfig.LikeResponse})
}

func (h LikeHandlers) dislikeHandler(c tb.Context) error {
	_ = c.Respond()
	// Находим сообщение по TelegramMessageID
	messageID := c.Message().ID
	message, err := services.FindByTelegramID(messageID)
	if err != nil {
		return c.Respond(&tb.CallbackResponse{Text: "Сообщение не найдено."})
	}

	err = services.SetDisliked(message)
	if err != nil {
		return c.Respond(&tb.CallbackResponse{Text: "Ошибка при обновлении оценки."})
	}

	// Отправляем уведомление в сервис ИИ
	aiClient := services.AIClient{}
	if err := aiClient.DislikeQuestion(*message.AIMessageID); err != nil {
		fmt.Printf("Ошибка при отправке оценки в ИИ сервис: %v\n", err)
	}

	services.ClearSourceMessage(c)

	return c.Respond(&tb.CallbackResponse{Text: h.likeConfig.DislikeResponse})
}
