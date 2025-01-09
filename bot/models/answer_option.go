package models

import (
	"github.com/google/uuid"
	"time"
)

type AnswerOption struct {
	ID               uuid.UUID       `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserId           *uuid.UUID      `gorm:"type:uuid;not null"`
	QuestionId       uuid.UUID       `gorm:"type:uuid;not null"`
	QuestionOptionId uuid.UUID       `gorm:"type:uuid;not null"`
	AnswerId         uuid.UUID       `gorm:"type:uuid;not null"`
	CreatedAt        time.Time       `gorm:"autoCreateTime"`              // Использует now() по умолчанию
	UpdatedAt        time.Time       `gorm:"autoUpdateTime"`              // Автоматически обновляется при изменении модели
	DeletedAt        *time.Time      `gorm:"index"`                       // Для soft delete
	User             *User           `gorm:"foreignKey:UserId"`           // Аннотация для указания внешнего ключа
	Question         *Question       `gorm:"foreignKey:QuestionId"`       // Аннотация для указания внешнего ключа
	QuestionOption   *QuestionOption `gorm:"foreignKey:QuestionOptionId"` // Аннотация для указания внешнего ключа
	Answer           *Answer         `gorm:"foreignKey:AnswerId"`         // Аннотация для указания внешнего ключа
}
