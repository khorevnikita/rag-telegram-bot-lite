package models

import (
	"github.com/google/uuid"
	"time"
)

type QuestionOption struct {
	ID                    uuid.UUID  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	QuestionId            uuid.UUID  `gorm:"type:uuid;not null"`
	Text                  string     `gorm:""`
	RequireAdditionalText bool       `gorm:"default:false"`
	CreatedAt             time.Time  `gorm:"autoCreateTime"`
	UpdatedAt             time.Time  `gorm:"autoUpdateTime"`
	DeletedAt             *time.Time `gorm:"index"`                 // Для soft delete
	Question              *Question  `gorm:"foreignKey:QuestionId"` // Аннотация для указания внешнего ключа
}
