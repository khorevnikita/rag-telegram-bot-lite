package services

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorag-telegram-bot/database"
	"gorag-telegram-bot/models"
	"gorm.io/gorm"
	"time"
)

func GetSubscription(user *models.User) (*models.Subscription, error) {
	var subscription models.Subscription
	result := database.DB.
		Where("user_id = ?", user.ID).
		Where("deleted_at is null").
		First(&subscription)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, nil
	}

	return &subscription, nil
}

func CreateSubscription(user *models.User) (*models.Subscription, error) {
	subscription := models.Subscription{
		UserID:    user.ID,
		Amount:    0,
		ExpiresAt: time.Now(),
	}
	if err := database.DB.Create(&subscription).Error; err != nil {
		fmt.Printf("Error saving message log: %v\n", err)
		return nil, err
	}
	return &subscription, nil
}

func GrantSubscription(user *models.User, method *models.PaymentMethod, amount float64) (*models.Subscription, error) {
	// Получение текущей подписки пользователя
	subscription, err := GetSubscription(user)
	if err != nil {
		return nil, err
	}

	// Если подписка отсутствует, создаем новую
	if subscription == nil {
		subscription, err = CreateSubscription(user)
		if err != nil {
			return nil, err
		}
	}

	// Устанавливаем дату окончания подписки на 1 месяц вперед
	now := time.Now()
	expireAt := now.AddDate(0, 1, 0) // 0 лет, 1 месяц, 0 дней

	// Обновляем данные подписки
	subscription.Amount = amount
	subscription.ExpiresAt = expireAt
	subscription.PaymentMethodID = &method.ID
	subscription.ProviderID = &method.ProviderID

	// Сохраняем обновленную подписку в базе данных
	result := database.DB.Save(subscription)
	if result.Error != nil {
		return nil, result.Error
	}

	return subscription, nil
}

func RenewSubscription(subscription *models.Subscription, time time.Time) error {
	subscription.ExpiresAt = time
	result := database.DB.Save(subscription)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func Unsubscribe(user *models.User) error {
	subscription, err := GetSubscription(user)
	if err != nil {
		return err
	}

	if subscription.UnsubscribedAt != nil {
		return nil //Уже отписался, сбивать дату не будем
	}

	now := time.Now()
	subscription.UnsubscribedAt = &now
	result := database.DB.Save(subscription)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

type SubscriptionReport struct {
	UserID                uuid.UUID  `json:"user_id"`
	TelegramID            int        `json:"telegram_id"`
	TelegramUsername      *string    `json:"telegram_username"`
	SubscriptionCreatedAt time.Time  `json:"subscription_created_at"`
	SubscriptionExpiresAt time.Time  `json:"subscription_expires_at"`
	UnsubscribedAt        *time.Time `json:"unsubscribed_at"`
	PaymentsTotalAmount   float64    `json:"payments_total_amount"`
	LastPaymentAt         *time.Time `json:"last_payment_at"`
}

func GetSubscriptionReport() ([]SubscriptionReport, error) {
	var reports []SubscriptionReport

	// Query to fetch the subscription reports
	err := database.DB.Table("subscriptions").
		Select(`users.id AS user_id, 
			users.telegram_id AS telegram_id, 
			users.telegram_username AS telegram_username, 
			subscriptions.created_at AS subscription_created_at, 
			subscriptions.expires_at AS subscription_expires_at, 
			subscriptions.unsubscribed_at AS unsubscribed_at, 
			COALESCE(SUM(payments.amount), 0) AS payments_total_amount, 
			MAX(payments.created_at) AS last_payment_at`).
		Joins("JOIN users ON subscriptions.user_id = users.id").
		Joins("LEFT JOIN payments ON subscriptions.id = payments.subscription_id").
		Group("subscriptions.id, users.id").
		Scan(&reports).Error

	if err != nil {
		return nil, err
	}

	return reports, nil
}
