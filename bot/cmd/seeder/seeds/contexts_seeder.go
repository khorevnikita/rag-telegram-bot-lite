package seeds

import (
	"fmt"
	"gorag-telegram-bot/config"
	"gorag-telegram-bot/models"
	"gorm.io/gorm"
	"time"
)

func SeedContexts(db *gorm.DB) {
	items := config.AppConfig.Menu.Items

	// Помечаем существующие записи как удалённые
	if err := markContextsDeleted(db); err != nil {
		fmt.Printf("Error marking deleted records: %v\n", err)
		return
	}

	// Добавляем вопросы в базу данных
	for _, i := range items {
		if i.Context == "" {
			continue
		}

		question := models.SystemContext{
			Key:  i.Key,
			Text: i.Context,
		}

		if err := db.Create(&question).Error; err != nil {
			fmt.Printf("Could not seed questions: %v\n", err)
			continue
		}
	}
}

func markContextsDeleted(db *gorm.DB) error {
	now := time.Now()
	if err := db.Model(&models.SystemContext{}).Where("deleted_at IS NULL").Update("deleted_at", now).Error; err != nil {
		return fmt.Errorf("Could not mark deleted in contexts: %v", err)
	}
	return nil
}
