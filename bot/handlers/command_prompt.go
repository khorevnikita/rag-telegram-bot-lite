package handlers

import (
	tb "gopkg.in/telebot.v4"
	"gorag-telegram-bot/config"
	"gorag-telegram-bot/core"
	"gorag-telegram-bot/services"
)

type CommandPromptHandlers struct {
	AIClient  services.AIClient
	Commands  []config.CommandConfig
	MenuItems []config.MenuItem
}

func GetNewCommandPromptHandlers() CommandPromptHandlers {
	return CommandPromptHandlers{
		AIClient:  services.NewAIClient(),
		Commands:  config.AppConfig.Commands,
		MenuItems: config.AppConfig.Menu.Items,
	}
}

func (h CommandPromptHandlers) RegisterCommands() {
	core.Bot.Handle(&tb.InlineButton{Unique: "command_prompt"}, func(c tb.Context) error {
		_ = c.Respond()
		user, err := services.GetUserByTelegramID(c.Sender().ID)
		if err != nil {
			return c.Send("Пользователь не найден. Пожалуйста, начните с команды /start.")
		}

		if !config.AppConfig.Modules.Form.CanSkip && user.FormCompletedAt == nil {
			return services.SendStartFormMessage(c.Sender().ID)
		}

		services.IncreaseUserMessagesStats(user)

		tempMsg, err := services.SendTemporaryMessage(c)
		if err != nil {
			return err
		}

		var prompt string
		for _, command := range h.Commands {
			for _, act := range command.Actions {
				if act.ActData == c.Data() {
					prompt = act.Prompt
					break
				}
			}
		}

		if prompt == "" {
			for _, menu := range h.MenuItems {
				for _, a := range menu.Actions {
					if a.ActData == c.Data() {
						prompt = a.Prompt
						break
					}
				}
			}
		}

		messageLog, err := services.SaveMessage(user, tempMsg, prompt, []services.QuestionFile{})
		if err != nil {
			return c.Send(err.Error())
		}

		_, err = h.AIClient.GetAIResponse(user, &messageLog, nil)
		if err != nil {
			return c.Send(err.Error())
		}

		return services.SendResponseMessage(c, tempMsg, &messageLog)
	})
}
