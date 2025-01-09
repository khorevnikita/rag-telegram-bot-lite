package services

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorag-telegram-bot/database"
	"gorag-telegram-bot/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

func GetQuestions() ([]models.Question, error) {
	var questions []models.Question

	// Основной запрос, который исключает вопросы, на которые уже есть ответы, и сортирует по полю `Order`
	result := database.DB.
		Preload("QuestionOptions").
		Where("questions.is_published = ?", true).
		Where("questions.deleted_at is null").
		Order("questions.order"). // Обновлено для сортировки в указанном порядке
		Find(&questions)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, nil
	}

	return questions, nil
}
func GetNextQuestion(user *models.User) (*models.Question, error) {
	var question models.Question
	// Подзапрос для идентификаторов вопросов, на которые пользователь уже ответил
	subQuery := database.DB.Model(&models.Answer{}).
		Select("question_id").
		Where("user_id = ?", user.ID).
		Where("answered_at is not null")

	// Основной запрос, который исключает вопросы, на которые уже есть ответы, и сортирует по полю `Order`
	result := database.DB.
		Preload("QuestionOptions").
		Where("questions.is_published = ?", true).
		Where("questions.deleted_at is null").
		Not("questions.id IN (?)", subQuery).
		Order("questions.order"). // Обновлено для сортировки в указанном порядке
		First(&question)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, nil
	}

	return &question, nil
}

func GetQuestionByID(qid *uuid.UUID) (*models.Question, error) {
	var question models.Question
	result := database.DB.
		Preload("QuestionOptions").
		Where("questions.is_published = ?", true).
		Where("questions.deleted_at is null").
		Where("questions.id = ?", qid).
		First(&question)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, nil
	}

	return &question, nil
}

func SaveAnswer(user *models.User, question *models.Question, text string, option *models.QuestionOption, completed bool) (*models.Answer, error) {
	var answeredAt *time.Time
	if completed {
		now := time.Now()
		answeredAt = &now
	}

	answer := models.Answer{
		UserId:     &user.ID,
		QuestionId: question.ID,
		Text:       text,
		AnsweredAt: answeredAt,
	}

	err := database.DB.Clauses(clause.Returning{}).Create(&answer).Error
	if err != nil {
		return nil, err
	}

	if option == nil {
		return &answer, nil
	} else {
		answerOpt := models.AnswerOption{
			UserId:           &user.ID,
			AnswerId:         answer.ID,
			QuestionId:       question.ID,
			QuestionOptionId: option.ID,
		}

		err = database.DB.Clauses(clause.Returning{}).Create(&answerOpt).Error
		if err != nil {
			return &answer, err
		}
		return &answer, nil
	}
}
func AddAnswerOption(answer *models.Answer, text string, option *models.QuestionOption) (*models.Answer, error) {
	newText := fmt.Sprintf("%s; %s", answer.Text, text)
	err := database.DB.Model(answer).Update("text", newText).Error
	if err != nil {
		return answer, err
	}

	answerOpt := models.AnswerOption{
		UserId:           answer.UserId,
		AnswerId:         answer.ID,
		QuestionId:       answer.QuestionId,
		QuestionOptionId: option.ID,
	}
	err = database.DB.Clauses(clause.Returning{}).Create(&answerOpt).Error
	if err != nil {
		return answer, err
	}
	return answer, nil
}

func SerializeUserAnswers(user *models.User) string {
	answers, err := GetUserAnswers(user)
	if err != nil {
		return ""
	}

	message := ""
	for _, answer := range answers {
		message += fmt.Sprintf("%s: *%s*\n", answer.Question.Text, answer.Text)
	}

	return message
}

func GetUserAnswers(user *models.User) ([]models.Answer, error) {
	var answers []models.Answer

	// Формирование подзапроса для получения последних дат ответов по каждому вопросу для пользователя
	subQuery := database.DB.Model(&models.Answer{}).
		Select("question_id, MAX(created_at) as last_answer_datetime").
		Where("user_id = ?", user.ID).
		Group("question_id")

	// Присоединяем подзапрос к основной таблице `answers` для получения полных записей этих последних ответов
	err := database.DB.Model(&models.Answer{}).
		Preload("Question").
		Select("answers.*").
		Joins("JOIN (?) as a on a.question_id = answers.question_id AND a.last_answer_datetime = answers.created_at", subQuery).
		Joins("JOIN questions as q on q.id = answers.question_id").
		Where("answers.user_id = ?", user.ID).
		Where("answers.text != ''").
		Order("q.order").
		Find(&answers).Error

	if err != nil {
		return nil, err
	}

	return answers, nil
}

func GetAnswers() ([]models.Answer, error) {
	var answers []models.Answer

	// Формирование подзапроса для получения последних дат ответов по каждому вопросу для каждого пользователя
	subQuery := database.DB.Model(&models.Answer{}).
		Select("user_id, question_id, MAX(created_at) as last_answer_datetime").
		Where("text != ''").
		Group("user_id, question_id")

	// Присоединяем подзапрос к основной таблице `answers` для получения полных записей этих последних ответов
	err := database.DB.Model(&models.Answer{}).
		//Preload("Question").
		//Preload("User").
		Select("answers.*").
		Joins("JOIN (?) as a on a.question_id = answers.question_id AND a.user_id = answers.user_id AND a.last_answer_datetime = answers.created_at", subQuery).
		//Joins("JOIN questions as q on q.id = answers.question_id").
		Where("answers.text != ''").
		Order("answers.user_id, answers.created_at desc").
		Find(&answers).Error

	if err != nil {
		return nil, err
	}

	return answers, nil
}

func GetQuestionAnswer(userID uuid.UUID, questionID uuid.UUID) (*models.Answer, error) {
	var answer models.Answer

	err := database.DB.Model(&models.Answer{}).
		Where("user_id = ?", userID).
		Where("question_id = ?", questionID).
		Order("created_at desc").
		First(&answer).Error

	if err != nil {
		// Если запись не найдена, вернуть nil вместо ошибки
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &answer, nil
}

func GetAnswerOptions(answerID uuid.UUID) ([]models.AnswerOption, error) {
	var options []models.AnswerOption

	err := database.DB.Model(&models.AnswerOption{}).
		Where("answer_id = ?", answerID).
		Find(&options).Error

	if err != nil {
		return nil, err
	}

	return options, nil
}

func MarkAsAnswered(answerID uuid.UUID) error {
	err := database.DB.Model(&models.Answer{}).
		Where("id = ?", answerID).
		Update("answered_at", "now()").Error

	if err != nil {
		return err
	}

	return nil
}

func FindAnswer(answerID uuid.UUID) (*models.Answer, error) {
	var answer models.Answer

	err := database.DB.Model(&models.Answer{}).
		Where("id = ?", answerID).
		Order("created_at desc").
		First(&answer).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &answer, nil
}
