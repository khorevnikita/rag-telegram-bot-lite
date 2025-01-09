package services

import (
	"errors"
	"github.com/google/uuid"
	"gorag-telegram-bot/database"
	"gorag-telegram-bot/models"
	"gorm.io/gorm"
)

func FindSystemContext(id *uuid.UUID) (*models.SystemContext, error) {
	var question models.SystemContext
	result := database.DB.
		Where("id = ?", id).
		Where("deleted_at is null").
		First(&question)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, nil
	}

	return &question, nil
}

func GetContextByKey(key string) (*models.SystemContext, error) {
	var question models.SystemContext
	result := database.DB.
		Where("key = ?", key).
		Where("deleted_at is null").
		First(&question)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, nil
	}

	return &question, nil
}
