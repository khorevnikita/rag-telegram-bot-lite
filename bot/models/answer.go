package models

import (
	"github.com/google/uuid"
	"time"
)

type Answer struct {
	ID            uuid.UUID  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserId        *uuid.UUID `gorm:"type:uuid"`
	QuestionId    uuid.UUID  `gorm:"type:uuid;not null"`
	Text          string     `gorm:""`
	CreatedAt     time.Time  `gorm:"autoCreateTime"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime"`
	AnsweredAt    *time.Time `gorm:"autoUpdateTime"`
	DeletedAt     *time.Time `gorm:"index"` // Для soft delete
	User          *User      `gorm:"foreignKey:UserId"`
	Question      *Question  `gorm:"foreignKey:QuestionId"`
	AnswerOptions []AnswerOption
}
