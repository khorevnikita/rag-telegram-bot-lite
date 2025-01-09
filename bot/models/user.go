package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID               uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TelegramID       int64     `gorm:"uniqueIndex"` // уникальный индекс
	TelegramUsername *string
	FirstName        *string
	LastName         *string
	Avatar           *string
	State            *string
	StateID          *uuid.UUID
	SystemContextID  *uuid.UUID `gorm:"index"` // Связь с SystemContext
	ConnectionDate   time.Time
	MessageCount     int
	LastMessageDate  time.Time
	ConversationID   *int
	IsAdmin          bool      `gorm:"default:false"` // Новое поле is_admin
	CreatedAt        time.Time `gorm:"autoCreateTime"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime"`
	FormCompletedAt  *time.Time
	Answers          []Answer
	Messages         []Message
	SystemContext    *SystemContext `gorm:"foreignKey:SystemContextID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // Внешний ключ
}

// Автоматически устанавливает UUID перед созданием записи
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
