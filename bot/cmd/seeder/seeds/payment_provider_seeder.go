package seeds

import (
	"fmt"
	"gorag-telegram-bot/config"
	"gorag-telegram-bot/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

func SeedPaymentProvider(db *gorm.DB) {
	billing := config.AppConfig.Modules.Billing

	// Помечаем существующие записи как удалённые
	if err := markPPsDeleted(db); err != nil {
		fmt.Printf("Error marking deleted records: %v\n", err)
		return
	}

	if billing.Providers.CloudPayments.Enabled {
		p := models.PaymentProvider{
			Name: models.ProviderCloudPayments,
		}
		if err := upsertPaymentProvider(db, &p); err != nil {
			fmt.Printf("Could not seed payment provider: %v\n", err)
		}
	}

	if billing.Providers.YooKassa.Enabled {
		p := models.PaymentProvider{
			Name: models.ProviderYooKassa,
		}
		if err := upsertPaymentProvider(db, &p); err != nil {
			fmt.Printf("Could not seed payment provider: %v\n", err)
		}
	}
}

func upsertPaymentProvider(db *gorm.DB, provider *models.PaymentProvider) error {
	// Используем INSERT ... ON CONFLICT для PostgreSQL
	err := db.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},                                // Уникальное поле для проверки
			DoUpdates: clause.AssignmentColumns([]string{"deleted_at", "updated_at"}), // Поля для обновления
		},
	).Create(provider).Error
	if err != nil {
		return fmt.Errorf("Could not upsert payment provider: %v", err)
	}
	return nil
}

func markPPsDeleted(db *gorm.DB) error {
	now := time.Now()
	if err := db.Model(&models.PaymentProvider{}).Where("deleted_at IS NULL").Update("deleted_at", now).Error; err != nil {
		return fmt.Errorf("Could not mark deleted in payment providers: %v", err)
	}
	return nil
}
