package handlers

import (
	"fmt"
	tb "gopkg.in/telebot.v4"
	"gorag-telegram-bot/config"
	"gorag-telegram-bot/core"
	"gorag-telegram-bot/services"
	"gorag-telegram-bot/utils"
	"log"
	"strings"
	"time"
)

type CommandHandlers struct {
}

func GetNewCommandHandler() CommandHandlers {
	return CommandHandlers{}
}

func (h CommandHandlers) RegisterCommands() {
	for _, cmd := range config.AppConfig.Commands {
		if cmd.Enabled {
			// Создаём локальную копию переменной cmd
			cmdCopy := cmd
			name := cmdCopy.Name

			switch name {
			case "start":
				core.Bot.Handle("/"+name, func(context tb.Context) error {
					return h.startHandler(context, cmdCopy)
				})
			case "menu":
				core.Bot.Handle("/"+name, func(context tb.Context) error {
					return h.menuHandler(context)
				})
			case "report":
				core.Bot.Handle("/"+name, func(context tb.Context) error {
					return h.reportHandler(context)
				})
			case "form":
				core.Bot.Handle("/"+name, func(context tb.Context) error {
					return h.formHandler(context)
				})
			case "subscription":
				core.Bot.Handle("/"+name, func(context tb.Context) error {
					return h.subscriptionHandler(context)
				})
			default:
				core.Bot.Handle("/"+name, func(context tb.Context) error {
					return h.handleRandomCommand(context, cmdCopy)
				})
			}
		}
	}
}

func (h CommandHandlers) startHandler(c tb.Context, cmd config.CommandConfig) error {
	user := c.Sender()
	err := services.SaveUserInfo(user, nil)
	if err != nil {
		log.Println("Failed to save user info:", err)
	}

	var inlineKeys [][]tb.InlineButton
	if len(cmd.Actions) > 0 {
		for _, act := range cmd.Actions {
			actBtn := tb.InlineButton{Unique: act.ActUnique, Text: act.Label, Data: act.ActData}
			inlineKeys = append(inlineKeys, []tb.InlineButton{actBtn})
		}
	}

	err = c.Send(cmd.Message, &tb.SendOptions{ParseMode: tb.ModeMarkdownV2}, &tb.ReplyMarkup{
		OneTimeKeyboard: true,
		InlineKeyboard:  inlineKeys,
	})
	if err != nil {
		return err
	}

	conf := config.AppConfig
	if conf.Modules.Form.Enabled && conf.Modules.Form.ShowOnStart {
		return services.SendStartFormMessage(c.Sender().ID)
	}

	return nil
}

// Обработчик стандартной команды пользователя
func (h CommandHandlers) handleRandomCommand(c tb.Context, cmd config.CommandConfig) error {
	var replyKeys [][]tb.ReplyButton

	if len(cmd.Replies) > 0 {
		for _, text := range cmd.Replies {
			optBtn := tb.ReplyButton{Text: text}
			replyKeys = append(replyKeys, []tb.ReplyButton{optBtn})
		}
	}

	var inlineKeys [][]tb.InlineButton
	if len(cmd.Actions) > 0 {
		for _, act := range cmd.Actions {
			actBtn := tb.InlineButton{Unique: act.ActUnique, Text: act.Label, Data: act.ActData}
			inlineKeys = append(inlineKeys, []tb.InlineButton{actBtn})
		}
	}

	return c.Send(cmd.Message, &tb.SendOptions{
		ParseMode: tb.ModeMarkdownV2,
	}, &tb.ReplyMarkup{
		ReplyKeyboard:   replyKeys,
		OneTimeKeyboard: true,
		InlineKeyboard:  inlineKeys,
	})
}

func (h CommandHandlers) reportHandler(c tb.Context) error {
	// Проверяем, является ли пользователь администратором
	user, err := services.GetUserByTelegramID(c.Sender().ID)
	if err != nil || !user.IsAdmin {
		return c.Send("Извините, у вас нет доступа к этой команде.")
	}

	var inlineKeys [][]tb.InlineButton
	var row []tb.InlineButton

	row = append(row, tb.InlineButton{
		Unique: "report_users",
		Text:   "Анкеты пользователей",
	})

	if config.AppConfig.Modules.Billing.Enabled {
		subscriptionsListBtb := tb.InlineButton{
			Unique: "report_subscriptions",
			Text:   "Все подписки",
		}
		row = append(row, subscriptionsListBtb)
	}

	inlineKeys = append(inlineKeys, row)

	if config.AppConfig.Modules.Likes.Enabled {
		messageStats := tb.InlineButton{
			Unique: "report_dislike_messages",
			Text:   "Плохие ответы",
		}
		inlineKeys = append(inlineKeys, []tb.InlineButton{messageStats})
	}

	now := time.Now()
	reportStart := now.Add(-1 * time.Hour * 24 * 30)
	regFile := h.getRegistrationsGraph(reportStart)
	subFile := h.sendSubscriptionsGraph(reportStart)

	// Формируем группу медиафайлов
	album := tb.Album{
		&tb.Photo{
			File:    regFile,
			Caption: "График регистраций за последний месяц",
		},
		&tb.Photo{
			File:    subFile,
			Caption: "График подписок за последний месяц",
		},
	}

	// Отправляем альбом
	_ = c.SendAlbum(album)
	report := services.GetCommonReport()
	return c.Send(utils.EscapeMarkdownV2WithHeaders(report), &tb.SendOptions{
		ParseMode: tb.ModeMarkdownV2,
	}, &tb.ReplyMarkup{
		InlineKeyboard: inlineKeys,
	})
}

func (h CommandHandlers) getRegistrationsGraph(startTime time.Time) tb.File {
	registrationActivity := services.GetUserRegistrationsTimeline(startTime)
	imagePath := services.DrawLine(startTime, registrationActivity, "Регистрации", "Дата", "Кол-во пользователей")
	//photo := &tb.Photo{File: tb.FromDisk(imagePath)}
	return tb.FromDisk(imagePath)
}

func (h CommandHandlers) sendSubscriptionsGraph(startTime time.Time) tb.File {
	subscriptionActivity := services.GetSubscriptionsTimeline(startTime)
	imagePath := services.DrawLine(startTime, subscriptionActivity, "Подписки", "Дата", "Кол-во активных подписок")
	//photo := &tb.Photo{File: tb.FromDisk(imagePath)}
	return tb.FromDisk(imagePath)
}

func (h CommandHandlers) formHandler(c tb.Context) error {
	formConfig := config.AppConfig.Modules.Form
	if !formConfig.Enabled {
		return c.Send("Модуль выключен")
	}

	user, err := services.GetUserByTelegramID(c.Sender().ID)
	if err != nil {
		return c.Send("Ошибка системы. Пожалуйста, обратитесь в тех. поддержку")
	}

	if user.FormCompletedAt != nil {
		return services.SendEditFormMessage(c)
	}

	return services.SendStartFormMessage(c.Sender().ID)
}

func (h CommandHandlers) subscriptionHandler(c tb.Context) error {
	billingConfig := config.AppConfig.Modules.Billing
	if !billingConfig.Enabled {
		return c.Send("Модуль выключен")
	}

	user, err := services.GetUserByTelegramID(c.Sender().ID)
	if err != nil {
		return c.Send("Ошибка системы. Пожалуйста, обратитесь в тех. поддержку")
	}

	subscription, err := services.GetSubscription(user)
	if err != nil {
		fmt.Printf("Erorr getting subscription: %s\n", err.Error())
		return c.Send("Ошибка системы. Пожалуйста, обратитесь в тех. поддержку")
	}

	if subscription == nil {
		return services.SendSubscribeMessage(c)
	}

	if !subscription.IsActive() {
		return services.SendSubscribeMessage(c)
	}

	msg := strings.Replace(billingConfig.SubscriptionMessage, "{expires_at}", subscription.ExpiresAt.Format("15:04 02\\.01\\.2006"), 1)

	inlineKeys := [][]tb.InlineButton{
		{
			tb.InlineButton{
				Unique: "billing_unsubscribe",
				Text:   billingConfig.UnsubscribeBtn,
			},
		},
	}

	return c.Send(msg, &tb.SendOptions{
		ParseMode: tb.ModeMarkdownV2,
	}, &tb.ReplyMarkup{
		InlineKeyboard: inlineKeys,
	})
}

func (h CommandHandlers) menuHandler(c tb.Context) error {
	conf := config.AppConfig.Menu
	if !conf.Enabled {
		return c.Send("Меню выключено")
	}

	return services.SendMenuMessage(c.Sender().ID)
}
