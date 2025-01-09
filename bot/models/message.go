package models

import (
	"github.com/google/uuid"
	"time"
)

type Message struct {
	ID                uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID            uuid.UUID
	TelegramMessageID int  // ID сообщения в Telegram
	AIMessageID       *int // ID сообщения от ИИ
	SystemContextID   *uuid.UUID
	Content           string
	Response          *string
	Liked             *bool          // Указатель на bool для хранения оценки (nil - оценки нет, true - лайк, false - дизлайк)
	CreatedAt         time.Time      `gorm:"autoCreateTime"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime"`
	User              *User          `gorm:"foreignKey:UserID"` // Аннотация для указания внешнего ключа
	SystemContext     *SystemContext `gorm:"foreignKey:SystemContextID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
