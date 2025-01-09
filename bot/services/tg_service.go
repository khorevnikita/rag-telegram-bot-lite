package services

import (
	"fmt"
	"github.com/google/uuid"
	tb "gopkg.in/telebot.v4"
	"gorag-telegram-bot/config"
	"gorag-telegram-bot/core"
	"gorag-telegram-bot/database"
	"gorag-telegram-bot/models"
	"gorag-telegram-bot/utils"
)

const maxTelegramMessageLength = 4096

func SendStartFormMessage(tgID int64) error {
	// Создаем кнопки
	formConfig := config.AppConfig.Modules.Form
	btnYes := tb.InlineButton{Unique: "form_start", Text: formConfig.StartLabel}
	// Создаем инлайн клавиатуру
	inlineKeys := [][]tb.InlineButton{
		{btnYes},
	}

	if formConfig.CanSkip {
		btnLater := tb.InlineButton{Unique: "form_later", Text: formConfig.LaterLabel}
		inlineKeys = append(inlineKeys, []tb.InlineButton{btnLater})
	}

	keyboard := tb.ReplyMarkup{
		InlineKeyboard: inlineKeys,
	}

	_, err := core.Bot.Send(&tb.User{ID: tgID}, formConfig.DisclaimerText, &tb.SendOptions{
		ParseMode:   tb.ModeMarkdownV2,
		ReplyMarkup: &keyboard,
	})

	return err
}

func SendEditFormMessage(c tb.Context) error {
	formConfig := config.AppConfig.Modules.Form
	viewBtn := tb.InlineButton{Unique: "view_form", Text: formConfig.ViewLabel}
	inlineKeys := [][]tb.InlineButton{
		{viewBtn},
	}

	if formConfig.AllowEdit {
		editBtn := tb.InlineButton{Unique: "edit_form", Text: formConfig.EditLabel}
		inlineKeys = append(inlineKeys, []tb.InlineButton{editBtn})
	}

	return c.Send(formConfig.CompletedMessage, &tb.SendOptions{
		ParseMode: tb.ModeMarkdownV2,
	}, &tb.ReplyMarkup{
		InlineKeyboard: inlineKeys,
	})
}

func SendNextQuestion(c tb.Context) error {
	user, err := GetUserByTelegramID(c.Sender().ID)
	if err != nil {
		fmt.Printf("Error getting user %s\n", err.Error())
		return c.Send("Непредвиденная ошибка. Если повторяется при повторном вызове, обратитесь в тех. поддержку.")
	}

	nextQuestion, err := GetNextQuestion(user)
	if err != nil {
		fmt.Printf("Error getting user %s\n", err.Error())
		return c.Send("Непредвиденная ошибка. Если повторяется при повторном вызове, обратитесь в тех. поддержку.")
	}

	if nextQuestion == nil {
		return ProcessFormCompleted(user)
	}

	SetUserState(user, utils.PointerToString("form"), &nextQuestion.ID)

	var replyKeys [][]tb.ReplyButton
	if nextQuestion.Type == models.QuestionTypeSelect {
		for _, option := range nextQuestion.QuestionOptions {
			optBtn := tb.ReplyButton{Text: option.Text}
			replyKeys = append(replyKeys, []tb.ReplyButton{optBtn})
		}
	}

	msg := utils.EscapeMarkdownV2WithHeaders(nextQuestion.Text)
	if nextQuestion.Hint != nil {
		hint := *nextQuestion.Hint
		msg += fmt.Sprintf("\n\n_%s_", utils.EscapeMarkdownV2WithHeaders(hint)) // Для курсива используем одинарные подчеркивания
	}

	return c.Send(msg, &tb.SendOptions{ParseMode: tb.ModeMarkdownV2}, &tb.ReplyMarkup{
		ReplyKeyboard:   replyKeys,
		OneTimeKeyboard: true,
		ForceReply:      true,
	})
}

func ProcessFormCompleted(user *models.User) error {
	SetUserState(user, nil, nil)
	if user.FormCompletedAt == nil {
		SetUserFormCompleted(user)

		_ = database.Publish("form_completed", struct {
			UserId uuid.UUID `json:"user_id"`
		}{
			UserId: user.ID,
		})
	}
	formConfig := config.AppConfig.Modules.Form
	// Создаем кнопки
	viewBtn := tb.InlineButton{Unique: "view_form", Text: formConfig.ViewLabel}
	inlineKeys := [][]tb.InlineButton{
		{viewBtn},
	}

	if formConfig.AllowEdit {
		editBtn := tb.InlineButton{Unique: "edit_form", Text: formConfig.EditLabel}
		inlineKeys = append(inlineKeys, []tb.InlineButton{editBtn})
	}

	keyboard := tb.ReplyMarkup{
		InlineKeyboard: inlineKeys,
	}

	_, _ = core.Bot.Send(&tb.User{ID: user.TelegramID}, formConfig.CompletedMessage, &tb.SendOptions{
		ParseMode:   tb.ModeMarkdownV2,
		ReplyMarkup: &keyboard,
	})

	return SendAIModeDisclaimer(user.TelegramID)
}

func SendAIModeDisclaimer(tgID int64) error {
	if config.AppConfig.Menu.Enabled {
		return SendMenuMessage(tgID)
	}
	_, err := core.Bot.Send(&tb.User{ID: tgID}, config.AppConfig.AIModeMessage, &tb.SendOptions{ParseMode: tb.ModeMarkdownV2})
	return err
}

func ClearSourceMessage(c tb.Context) {
	// Убираем кнопки из сообщения
	var keyboard *tb.ReplyMarkup
	if config.AppConfig.Menu.Enabled {
		// Добавляем новый ряд кнопок в InlineKeyboard
		newRow := []tb.InlineButton{{
			Unique: "menu_button",
			Text:   config.AppConfig.Menu.Label, // Например, текст из конфигурации
		}}
		keyboard = &tb.ReplyMarkup{InlineKeyboard: [][]tb.InlineButton{newRow}}
	}

	_, _ = c.Bot().Edit(c.Callback().Message, utils.EscapeMarkdownV2WithHeaders(c.Callback().Message.Text), &tb.SendOptions{
		ParseMode:   tb.ModeMarkdownV2,
		ReplyMarkup: keyboard, // Убираем кнопки
	})
}

func SendTemporaryMessage(c tb.Context) (*tb.Message, error) {
	tempMsg, err := c.Bot().Send(c.Sender(), config.AppConfig.TemporaryMessage, &tb.SendOptions{
		ReplyMarkup: nil,
	})
	if err != nil {
		return nil, c.Send("Ошибка отправки промежуточного сообщения.")
	}
	return tempMsg, nil
}

func SendResponseMessage(c tb.Context, tempMsg *tb.Message, msgLog *models.Message) error {
	var keyboard tb.ReplyMarkup

	likesModule := config.AppConfig.Modules.Likes
	if likesModule.Enabled {
		// Создаем кнопки "Лайк" и "Дисклайк"
		likeButton := tb.InlineButton{
			Unique: "like_button",
			Text:   likesModule.LikeLabel,
		}
		dislikeButton := tb.InlineButton{
			Unique: "dislike_button",
			Text:   likesModule.DislikeLabel,
		}
		keyboard = tb.ReplyMarkup{
			InlineKeyboard: [][]tb.InlineButton{
				{likeButton, dislikeButton},
			},
		}
	}

	if config.AppConfig.Menu.Enabled {
		// Добавляем новый ряд кнопок в InlineKeyboard
		newRow := []tb.InlineButton{
			tb.InlineButton{
				Unique: "menu_button",
				Text:   config.AppConfig.Menu.Label, // Например, текст из конфигурации
			},
		}
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, newRow)
	}

	text := utils.EscapeMarkdownV2WithHeaders(*msgLog.Response)

	// Если текст короче лимита, отправляем как есть
	if len(text) <= maxTelegramMessageLength {
		_, err := c.Bot().Edit(tempMsg, text, &tb.SendOptions{
			ParseMode:   tb.ModeMarkdownV2,
			ReplyMarkup: &keyboard,
		})
		return err
	}

	// Разбиваем текст на части
	parts := utils.SplitTextIntoChunks(text, maxTelegramMessageLength)
	// Отправляем первую часть, редактируя tempMsg
	_, err := c.Bot().Edit(tempMsg, parts[0], &tb.SendOptions{
		ParseMode:   tb.ModeMarkdownV2,
		ReplyMarkup: nil,
	})
	if err != nil {
		return err
	}

	var lastMessage *tb.Message
	// Отправляем оставшиеся части как новые сообщения
	for k, part := range parts[1:] {
		var withKeyboard *tb.ReplyMarkup
		if k == (len(parts[1:]) - 1) {
			withKeyboard = &keyboard
		}

		lastMessage, err = c.Bot().Send(c.Recipient(), part, &tb.SendOptions{
			ParseMode:   tb.ModeMarkdownV2,
			ReplyMarkup: withKeyboard,
		})
		if err != nil {
			return err
		}
	}

	if lastMessage != nil {
		ChangeTelegramId(msgLog, lastMessage.ID)
	}
	return nil
}

func SendSubscribeMessage(c tb.Context) error {
	user, _ := GetUserByTelegramID(c.Sender().ID)

	billingConf := config.AppConfig.Modules.Billing

	if billingConf.Providers.YooKassa.Enabled {
		payButton := tb.Invoice{
			Title:       "Подписка на DIMA",
			Description: "1 мес. использования",
			Payload:     user.ID.String(),
			Token:       billingConf.Providers.YooKassa.Token, // Токен провайдера оплаты
			Currency:    "RUB",
			Prices: []tb.Price{
				{Label: "Подписка", Amount: billingConf.Providers.YooKassa.Price}, // Цена указывается в центах (например, 10.00 USD)
			},
		}

		// Отправка инвойса
		return c.Send(&payButton)
	}

	if billingConf.Providers.CloudPayments.Enabled {
		redirectURL := fmt.Sprintf("%s/api/checkout/cloud-payments?user_id=%s", config.AppConfig.AppURL, user.ID)

		payBtn := tb.InlineButton{
			URL:  redirectURL,
			Text: config.AppConfig.Modules.Billing.SubscribeBtn,
		}

		keyboard := tb.ReplyMarkup{
			InlineKeyboard: [][]tb.InlineButton{
				{payBtn},
			},
		}

		return c.Send(config.AppConfig.Modules.Billing.SubscriptionAlert, &tb.SendOptions{
			ParseMode:             tb.ModeMarkdownV2,
			ReplyMarkup:           &keyboard,
			DisableWebPagePreview: true,
		})
	}

	return c.Send("Приём платежей не подключен")
}

func SendNotification(user *models.User, text string) error {
	_, err := core.Bot.Send(&tb.User{ID: user.TelegramID}, text, config.AppConfig.Modules.Billing.SubscriptionAlert, &tb.SendOptions{
		ParseMode:             tb.ModeMarkdownV2,
		DisableWebPagePreview: true,
	})
	return err
}

func SendMenuMessage(tgID int64) error {
	menu := config.AppConfig.Menu
	var inlineKeys [][]tb.InlineButton

	var row []tb.InlineButton
	for i, m := range menu.Items {
		if !m.Enabled {
			continue
		}
		messageStats := tb.InlineButton{
			Unique: "menu_item",
			Text:   m.ButtonLabel,
			Data:   m.Key,
		}
		row = append(row, messageStats)

		if (i+1)%2 == 0 {
			inlineKeys = append(inlineKeys, row)
			row = nil // Очищаем ряд
		}
	}

	if row != nil {
		//Если нечетное кол-во меню
		inlineKeys = append(inlineKeys, row)
	}

	commands := config.AppConfig.Commands

	var msg string
	for _, com := range commands {
		if com.Name == "menu" {
			msg = com.Message
		}
	}

	keyboard := tb.ReplyMarkup{
		InlineKeyboard: inlineKeys,
	}

	_, err := core.Bot.Send(&tb.User{ID: tgID}, msg, &tb.SendOptions{
		ParseMode:   tb.ModeMarkdownV2,
		ReplyMarkup: &keyboard,
	})

	return err
}
