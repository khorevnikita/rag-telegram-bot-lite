package models

import (
	"github.com/google/uuid"
	"time"
)

type ProviderName string

const (
	ProviderCloudPayments ProviderName = "cloud_payments"
	ProviderYooKassa      ProviderName = "yoo_kassa"
)

type PaymentProvider struct {
	ID        uuid.UUID    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"` // Уникальный идентификатор
	Name      ProviderName `gorm:"type:varchar(50);unique;not null"`                // Название провайдера
	CreatedAt time.Time    `gorm:"default:now()"`                                   // Время создания
	UpdatedAt time.Time    `gorm:"default:now()"`                                   // Время обновления
	DeletedAt *time.Time   `gorm:"index"`                                           // Время удаления (для soft delete)
}
