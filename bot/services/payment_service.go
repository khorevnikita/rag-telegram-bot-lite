package services

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorag-telegram-bot/database"
	"gorag-telegram-bot/models"
	"gorag-telegram-bot/utils"
	"gorm.io/gorm"
	"strings"
)

func GetPaymentProvider(key models.ProviderName) (*models.PaymentProvider, error) {
	var provider models.PaymentProvider
	result := database.DB.
		Where("name = ?", key).
		Where("deleted_at is null").
		First(&provider)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, nil
	}

	return &provider, nil
}
func GetPaymentMethod(uid uuid.UUID, providerID uuid.UUID) (*models.PaymentMethod, error) {
	var method models.PaymentMethod
	result := database.DB.
		Where("user_id = ?", uid).
		Where("provider_id = ?", providerID).
		Where("deleted_at is null").
		First(&method)

	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, nil
	}

	return &method, nil
}

func SavePaymentMethod(uid uuid.UUID, providerID uuid.UUID, token string, lastFour string, extDate string) (*models.PaymentMethod, error) {
	existingMethod, err := GetPaymentMethod(uid, providerID)
	if err != nil {
		return nil, err
	}

	if existingMethod != nil {
		return UpdatePaymentMethod(existingMethod, token, lastFour, extDate)
	}
	return CreatePaymentMethod(uid, providerID, token, lastFour, extDate)
}
func CreatePaymentMethod(uid uuid.UUID, providerID uuid.UUID, token string, lastFour string, extDate string) (*models.PaymentMethod, error) {
	parts := strings.Split(extDate, "/")
	method := models.PaymentMethod{
		UserID:          uid,
		ProviderID:      providerID,
		Token:           token,
		CardLastFour:    lastFour,
		CardExpiryMonth: utils.ParseInt(parts[0]),
		CardExpiryYear:  utils.ParseInt(parts[1]),
	}
	if err := database.DB.Create(&method).Error; err != nil {
		fmt.Printf("Error saving message log: %v\n", err)
		return nil, err
	}

	return &method, nil
}
func UpdatePaymentMethod(method *models.PaymentMethod, token string, lastFour string, extDate string) (*models.PaymentMethod, error) {
	parts := strings.Split(extDate, "/")

	method.Token = token
	method.CardLastFour = lastFour
	method.CardExpiryMonth = utils.ParseInt(parts[0])
	method.CardExpiryYear = utils.ParseInt(parts[1])

	err := database.DB.Save(method).Error
	return method, err
}
func SavePayment(subID uuid.UUID, providerID uuid.UUID, amount float64) (*models.Payment, error) {
	payment := models.Payment{
		SubscriptionID: subID,
		ProviderID:     providerID,
		Amount:         amount,
		Currency:       "RUB",
	}
	if err := database.DB.Create(&payment).Error; err != nil {
		fmt.Printf("Error saving message log: %v\n", err)
		return nil, err
	}

	return &payment, nil
}
