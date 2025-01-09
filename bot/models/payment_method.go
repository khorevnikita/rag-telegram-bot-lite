package models

import (
	"github.com/google/uuid"
	"time"
)

type PaymentMethod struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"` // Уникальный идентификатор
	UserID          uuid.UUID  `gorm:"type:uuid;not null"`                              // Ссылка на пользователя
	ProviderID      uuid.UUID  `gorm:"type:uuid;not null"`                              // Ссылка на провайдера
	Token           string     `gorm:"type:varchar(255);not null"`                      // Токен метода оплаты
	CardLastFour    string     `gorm:"type:varchar(4)"`                                 // Последние 4 цифры карты
	CardExpiryMonth int        `gorm:"type:int"`                                        // Месяц окончания действия карты
	CardExpiryYear  int        `gorm:"type:int"`                                        // Год окончания действия карты
	CreatedAt       time.Time  `gorm:"default:now()"`                                   // Время создания
	UpdatedAt       time.Time  `gorm:"default:now()"`                                   // Время обновления
	DeletedAt       *time.Time `gorm:"index"`                                           // Время удаления (для soft delete)

	// Relations
	User     User            `gorm:"foreignKey:UserID"`     // Связь с пользователем
	Provider PaymentProvider `gorm:"foreignKey:ProviderID"` // Связь с провайдером
}
