package handlers

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	tb "gopkg.in/telebot.v4"
	"gorag-telegram-bot/config"
	"gorag-telegram-bot/core"
	"gorag-telegram-bot/models"
	"gorag-telegram-bot/services"
	"gorag-telegram-bot/utils"
	"regexp"
	"strconv"
)

type FormHandlers struct {
	formConfig config.FormModule
}

func GetNewFormHandlers() FormHandlers {
	formConfig := config.AppConfig.Modules.Form
	return FormHandlers{
		formConfig: formConfig,
	}
}

func (h FormHandlers) RegisterCommands() {
	if h.formConfig.Enabled {
		core.Bot.Handle(&tb.InlineButton{Unique: "form_start"}, h.formStartHandler)
		if h.formConfig.CanSkip {
			core.Bot.Handle(&tb.InlineButton{Unique: "form_later"}, h.laterHandler)
		}

		core.Bot.Handle(&tb.InlineButton{Unique: "view_form"}, h.viewHandler)
		core.Bot.Handle(&tb.InlineButton{Unique: "form_answer_more"}, h.continueAnswer)
		core.Bot.Handle(&tb.InlineButton{Unique: "form_answer_complete"}, h.completeAnswer)

		if h.formConfig.AllowEdit {
			core.Bot.Handle(&tb.InlineButton{Unique: "edit_form"}, h.editHandler)
			core.Bot.Handle(&tb.InlineButton{Unique: "edit_question"}, h.editQuestionHandler)
		}
	}
}

func (h FormHandlers) formStartHandler(c tb.Context) error {
	_ = c.Respond()
	services.ClearSourceMessage(c)
	return services.SendNextQuestion(c)
}

func (h FormHandlers) laterHandler(c tb.Context) error {
	_ = c.Respond()
	services.ClearSourceMessage(c)
	_ = c.Send(h.formConfig.LaterMessage)
	return services.SendAIModeDisclaimer(c.Sender().ID)
}

func (h FormHandlers) viewHandler(c tb.Context) error {
	_ = c.Respond()
	services.ClearSourceMessage(c)

	user, err := services.GetUserByTelegramID(c.Sender().ID)
	if err != nil {
		fmt.Printf("Error getting user %s\n", err.Error())
		return c.Send("Непредвиденная ошибка. Если повторяется при повторном вызове, обратитесь в тех. поддержку.")
	}
	// Запрашиваем анкету пользователя в режиме чтения
	msg := services.SerializeUserAnswers(user)
	if msg == "" {
		return c.Send("Данные анкеты не найдены")
	}
	return c.Send(utils.EscapeMarkdownV2WithHeaders(msg), &tb.SendOptions{ParseMode: tb.ModeMarkdownV2})
}

func (h FormHandlers) editHandler(c tb.Context) error {
	_ = c.Respond()
	services.ClearSourceMessage(c)

	// 1. Вывести кнопками список вопросов, который надо отредактировать
	questions, err := services.GetQuestions()

	if err != nil {
		return c.Send("Ошибка при получении списка вопросов.")
	}

	var inlineKeys [][]tb.InlineButton

	for _, question := range questions {
		questionBtn := tb.InlineButton{Unique: "edit_question", Text: question.Text, Data: question.ID.String()}
		inlineKeys = append(inlineKeys, []tb.InlineButton{questionBtn})
	}

	return c.Send(config.AppConfig.Modules.Form.SelectQuestionMessage, &tb.ReplyMarkup{
		InlineKeyboard: inlineKeys,
		ForceReply:     true,
	})

}

func (h FormHandlers) editQuestionHandler(c tb.Context) error {
	_ = c.Respond()
	user, err := services.GetUserByTelegramID(c.Sender().ID)
	if err != nil {
		return c.Send("Непредвиденная ошибка")
	}

	services.ClearSourceMessage(c)

	questionID, err := uuid.Parse(c.Data())
	if err != nil {
		return c.Send("Не можем распознать вопрос")
	}

	// 1. Вывести кнопками список вопросов, который надо отредактировать
	question, err := services.GetQuestionByID(&questionID)

	if err != nil {
		return c.Send("Ошибка при получении вопроса")
	}

	services.SetUserState(user, utils.PointerToString("form"), &question.ID)
	// Запросить все вопросы: сколько отвечено, текущее неотвеченное.

	var replyKeys [][]tb.ReplyButton
	if question.Type == models.QuestionTypeSelect {
		for _, option := range question.QuestionOptions {
			optBtn := tb.ReplyButton{Text: option.Text}
			replyKeys = append(replyKeys, []tb.ReplyButton{optBtn})
		}
	}

	return c.Send(question.Text, &tb.ReplyMarkup{
		ReplyKeyboard:   replyKeys,
		OneTimeKeyboard: true,
		ForceReply:      true,
	})
}

func (h FormHandlers) continueAnswer(c tb.Context) error {
	_ = c.Respond()
	services.ClearSourceMessage(c)

	answerID, err := uuid.Parse(c.Data())
	if err != nil {
		return c.Send(err.Error())
	}
	answer, err := services.FindAnswer(answerID)
	if err != nil {
		return c.Send(err.Error())
	}
	question, err := services.GetQuestionByID(&answer.QuestionId)
	if err != nil {
		return c.Send(err.Error())
	}

	answeredOptions, err := services.GetAnswerOptions(answerID)
	if err != nil {
		return c.Send(err.Error())
	}

	var replyKeys [][]tb.ReplyButton
	for _, questionOption := range question.QuestionOptions {
		needToAdd := true
		for _, answerOption := range answeredOptions {
			if questionOption.ID == answerOption.QuestionOptionId {
				// Не добавляем
				needToAdd = false
			}
		}
		if needToAdd {
			optionToAdd := questionOption
			optBtn := tb.ReplyButton{Text: optionToAdd.Text}
			replyKeys = append(replyKeys, []tb.ReplyButton{optBtn})
		}
	}

	return c.Send(config.AppConfig.Modules.Form.AddOptionMessage, &tb.ReplyMarkup{
		ReplyKeyboard:   replyKeys,
		OneTimeKeyboard: true,
		ForceReply:      true,
	})
}

func (h FormHandlers) completeAnswer(c tb.Context) error {
	_ = c.Respond()
	services.ClearSourceMessage(c)

	answerID, err := uuid.Parse(c.Data())
	if err != nil {
		return c.Send("Не можем распознать ответ")
	}

	err = services.MarkAsAnswered(answerID)

	if err != nil {
		return c.Send("Ошибка при сохранении")
	}

	return services.SendNextQuestion(c)
}

func (h FormHandlers) ProcessFormAnswer(c tb.Context, user *models.User) error {
	var err error
	sendNext := true

	question, err := services.GetQuestionByID(user.StateID)
	if err != nil {
		return err
	}

	if question == nil {
		return fmt.Errorf("вопрос не найден")
	}

	switch question.Type {
	case models.QuestionTypeSelect:
		err, sendNext = h.processSelectAnswer(c, user, question)
	case models.QuestionTypeEmail:
		err = h.processEmailAnswer(c, user, question)
	case models.QuestionTypeNumber:
		err = h.processNumberAnswer(c, user, question)
	case models.QuestionTypeText:
		err = h.processTextAnswer(c, user, question)
	default:
		fmt.Printf("NEVER SHOULD BEEN HERE ERROR\n")
		err = fmt.Errorf("неизвестный тип вопроса")
	}

	if err != nil {
		_ = c.Send(err.Error())
	}

	if !sendNext {
		return nil
	}

	return services.SendNextQuestion(c)
}

func (h FormHandlers) processSelectAnswer(c tb.Context, user *models.User, question *models.Question) (error, bool) {
	freeOptionAvailable := false

	for _, option := range question.QuestionOptions {
		if option.RequireAdditionalText {
			freeOptionAvailable = true
		}

		if option.Text == c.Text() {
			if option.RequireAdditionalText {
				return c.Send(config.AppConfig.Modules.Form.CustomOptionMessage, &tb.ReplyMarkup{
					RemoveKeyboard: true,
				}), false
			} else {
				isOnlyOneOption := question.SelectableOptionsCount == 1
				fmt.Printf("isOnlyOneOption: %v\n", isOnlyOneOption)
				if isOnlyOneOption {
					_, err := services.SaveAnswer(user, question, c.Text(), &option, true)
					if err != nil {
						return err, false
					}
					return nil, true
				}

				answer, err := services.GetQuestionAnswer(user.ID, question.ID)
				if err != nil {
					return err, false
				}
				spew.Dump(answer)

				if answer == nil {
					answer, err = services.SaveAnswer(user, question, c.Text(), &option, false)
					if err != nil {
						return err, false
					}
				} else {
					// answer add option
					_, err = services.AddAnswerOption(answer, c.Text(), &option)
					if err != nil {
						return err, false
					}
				}

				answeredOptions, err := services.GetAnswerOptions(answer.ID)
				if err != nil {
					return err, true
				}

				if len(answeredOptions) >= question.SelectableOptionsCount {
					err = services.MarkAsAnswered(answer.ID)
					return err, true
				}

				markAnswerBtn := tb.InlineButton{Unique: "form_answer_complete", Text: config.AppConfig.Modules.Form.NextQuestionLabel, Data: answer.ID.String()}
				moreOptionBtn := tb.InlineButton{Unique: "form_answer_more", Text: config.AppConfig.Modules.Form.MoreOptionLabel, Data: answer.ID.String()}
				inlineKeys := [][]tb.InlineButton{{markAnswerBtn}, {moreOptionBtn}}

				return c.Send(config.AppConfig.Modules.Form.OptionSavedMessage, &tb.ReplyMarkup{
					InlineKeyboard:  inlineKeys,
					OneTimeKeyboard: true,
					ForceReply:      true,
				}), false
			}
		}
	}

	if freeOptionAvailable {
		ans, err := services.SaveAnswer(user, question, c.Text(), nil, true)
		spew.Dump(freeOptionAvailable, ans)
		return err, true
	}

	return fmt.Errorf("%s", config.AppConfig.Modules.Form.WrongOptionMessage), true
}

func (h FormHandlers) processEmailAnswer(c tb.Context, user *models.User, question *models.Question) error {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	isEmailValid := regexp.MustCompile(emailRegex).MatchString(c.Text())

	if !isEmailValid {
		return fmt.Errorf("Введи корректный email")
	}
	_, err := services.SaveAnswer(user, question, c.Text(), nil, true)
	return err
}

func (h FormHandlers) processNumberAnswer(c tb.Context, user *models.User, question *models.Question) error {
	if _, err := strconv.Atoi(c.Text()); err != nil {
		return fmt.Errorf("Введите корректное число")
	}
	_, err := services.SaveAnswer(user, question, c.Text(), nil, true)
	return err
}

func (h FormHandlers) processTextAnswer(c tb.Context, user *models.User, question *models.Question) error {
	_, err := services.SaveAnswer(user, question, c.Text(), nil, true)
	return err
}
