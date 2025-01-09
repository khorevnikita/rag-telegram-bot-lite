package handlers

import (
	tb "gopkg.in/telebot.v4"
	"gorag-telegram-bot/config"
	"gorag-telegram-bot/core"
	"gorag-telegram-bot/services"
)

type MenuHandlers struct {
	Bot    *tb.Bot
	Config config.MenuConfig
}

func GetMenuHandlers() MenuHandlers {
	return MenuHandlers{
		Bot:    core.Bot,
		Config: config.AppConfig.Menu,
	}
}

func (h MenuHandlers) RegisterCommands() {
	if h.Config.Enabled {
		// add some handlers
		h.Bot.Handle(&tb.InlineButton{Unique: "menu_button"}, h.onMainMenu)
		h.Bot.Handle(&tb.InlineButton{Unique: "menu_item"}, h.onMenuClick)
	}
}

func (h MenuHandlers) onMainMenu(c tb.Context) error {
	_ = c.Respond()
	return services.SendMenuMessage(c.Sender().ID)
}

func (h MenuHandlers) onMenuClick(c tb.Context) error {
	_ = c.Respond()

	var menu config.MenuItem
	for _, m := range h.Config.Items {
		if m.Key == c.Data() {
			menu = m
			break
		}
	}
	if !menu.Enabled {
		return c.RespondText("Меню не найдено")
	}

	var inlineKeys [][]tb.InlineButton
	if len(menu.Actions) > 0 {
		var row []tb.InlineButton
		for i, act := range menu.Actions {
			actBtn := tb.InlineButton{Unique: act.ActUnique, Text: act.Label, Data: act.ActData}
			row = append(row, actBtn)

			// Когда в ряду 2 кнопки, добавляем его в массив и начинаем новый ряд
			if (i+1)%2 == 0 {
				inlineKeys = append(inlineKeys, row)
				row = nil // Очищаем ряд
			}
		}

		// Добавляем оставшиеся кнопки, если их количество нечетное
		if len(row) > 0 {
			inlineKeys = append(inlineKeys, row)
		}
	}

	sysContext, err := services.GetContextByKey(menu.Key)
	if err != nil {
		return c.Send(err.Error())
	}

	user, err := services.GetUserByTelegramID(c.Sender().ID)
	if err != nil {
		return c.Send(err.Error())
	}
	services.SetSystemContext(user, sysContext)

	return c.Send(menu.Message, &tb.SendOptions{
		ParseMode: tb.ModeMarkdownV2,
	}, &tb.ReplyMarkup{
		InlineKeyboard: inlineKeys,
	})
}
