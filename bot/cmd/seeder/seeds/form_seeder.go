package seeds

import (
	"fmt"
	"gorag-telegram-bot/config"
	"gorag-telegram-bot/models"
	"gorm.io/gorm"
	"time"
)

func SeedForm(db *gorm.DB) {
	questions := config.AppConfig.Modules.Form.Questions

	// Помечаем существующие записи как удалённые
	if err := markQuestionsDeleted(db); err != nil {
		fmt.Printf("Error marking deleted records: %v\n", err)
		return
	}

	// Добавляем вопросы в базу данных
	for _, q := range questions {
		question := models.Question{
			Text:                   q.Text,
			Order:                  q.Order,
			IsRequired:             q.IsRequired,
			IsPublished:            true,
			Type:                   q.Type,
			SelectableOptionsCount: q.SelectableOptionsCount,
			Hint:                   q.Hint,
		}

		if err := db.Omit("QuestionOptions").Create(&question).Error; err != nil {
			fmt.Printf("Could not seed questions: %v\n", err)
			continue
		}

		// Добавляем опции вопроса
		for _, opt := range q.Options {
			option := models.QuestionOption{
				Text:                  opt.Text,
				RequireAdditionalText: opt.RequireAdditionalText,
				QuestionId:            question.ID,
			}

			if err := db.Create(&option).Error; err != nil {
				fmt.Printf("Could not seed question_options: %v\n", err)
			}
		}
	}
}

func markQuestionsDeleted(db *gorm.DB) error {
	now := time.Now()

	// Обновляем deleted_at для всех вопросов
	if err := db.Model(&models.Question{}).Where("deleted_at IS NULL").Update("deleted_at", now).Error; err != nil {
		return fmt.Errorf("Could not mark deleted in questions: %v", err)
	}

	// Обновляем deleted_at для всех опций вопросов
	if err := db.Model(&models.QuestionOption{}).Where("deleted_at IS NULL").Update("deleted_at", now).Error; err != nil {
		return fmt.Errorf("Could not mark deleted in question_options: %v", err)
	}
	return nil
}
