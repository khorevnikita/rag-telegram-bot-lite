package models

import (
	"github.com/google/uuid"
	"time"
)

type Payment struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"` // Уникальный идентификатор
	SubscriptionID uuid.UUID `gorm:"type:uuid;not null"`                              // Ссылка на подписку
	ProviderID     uuid.UUID `gorm:"type:uuid;not null"`                              // Ссылка на провайдера
	Amount         float64   `gorm:"type:decimal(10,2);not null"`                     // Сумма платежа
	Currency       string    `gorm:"type:varchar(10);default:'RUB'"`                  // Валюта
	CreatedAt      time.Time `gorm:"default:now()"`                                   // Время создания
	UpdatedAt      time.Time `gorm:"default:now()"`                                   // Время обновления

	// Relations
	Subscription Subscription    `gorm:"foreignKey:SubscriptionID"` // Связь с подпиской
	Provider     PaymentProvider `gorm:"foreignKey:ProviderID"`     // Связь с провайдером
}
