package models

import (
	"github.com/google/uuid"
	"time"
)

type MessageFile struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	MessageID uuid.UUID
	UserID    uuid.UUID

	FilePath  string
	FileName  string
	Size      int
	Extension string
	FileType  string

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	User      *User     `gorm:"foreignKey:UserID"` // Аннотация для указания внешнего ключа
	Message   *Message  `gorm:"foreignKey:MessageID"`
}
