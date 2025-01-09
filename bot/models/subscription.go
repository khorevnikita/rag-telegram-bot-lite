package models

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID              uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID          uuid.UUID  `gorm:"type:uuid"`
	PaymentMethodID *uuid.UUID `gorm:"type:uuid"` // Ссылка на метод оплаты
	ProviderID      *uuid.UUID `gorm:"type:uuid"`
	Amount          float64
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
	ExpiresAt       time.Time
	UnsubscribedAt  *time.Time
	DeletedAt       *time.Time `gorm:"index"`
	User            User
	PaymentMethod   *PaymentMethod   `gorm:"foreignKey:PaymentMethodID"`
	Provider        *PaymentProvider `gorm:"foreignKey:ProviderID"`
}

func (m Subscription) IsActive() bool {
	now := time.Now()
	return m.ExpiresAt.After(now)
}
