package models

import (
	"github.com/google/uuid"
	"time"
)

type SystemContext struct {
	ID        uuid.UUID  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Key       string     `gorm:""`
	Text      string     `gorm:""`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
	DeletedAt *time.Time `gorm:"index"`                                                                     // Для soft delete
	Users     []User     `gorm:"foreignKey:SystemContextID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // Указание внешнего ключа
	Messages  []Message
}
