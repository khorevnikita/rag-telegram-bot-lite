package services

import (
	"github.com/google/uuid"
	tb "gopkg.in/telebot.v4"
	"gorag-telegram-bot/database"
	"gorag-telegram-bot/models"
	"gorm.io/gorm/clause"
	"time"
)

func FindUser(uid uuid.UUID) (*models.User, error) {
	var user models.User
	err := database.DB.Where("id = ?", uid).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByTelegramID(telegramID int64) (*models.User, error) {
	var user models.User
	err := database.DB.Where("telegram_id = ?", telegramID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func SaveUserInfo(user *tb.User, conversationID *int) error {
	u := models.User{
		TelegramID:       user.ID,
		TelegramUsername: &user.Username,
		FirstName:        &user.FirstName,
		LastName:         &user.LastName,
		ConnectionDate:   time.Now(),
		MessageCount:     1,
		LastMessageDate:  time.Now(),
		ConversationID:   conversationID,
	}
	return database.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&u).Error
}

func SetUserFormCompleted(user *models.User) {
	now := time.Now()
	user.FormCompletedAt = &now
	database.DB.Save(user)
}

func SetUserConversation(user *models.User, conversationID *int) {
	user.ConversationID = conversationID
	database.DB.Save(user)
}

func SetUserState(user *models.User, state *string, stateID *uuid.UUID) {
	user.State = state
	user.StateID = stateID
	database.DB.Save(user)
}

func SetSystemContext(user *models.User, sysContext *models.SystemContext) {
	if sysContext != nil {
		user.SystemContextID = &sysContext.ID
	} else {
		user.SystemContextID = nil
	}

	database.DB.Save(user)
}

func GetUsers() ([]models.User, error) {
	var users []models.User
	if err := database.DB.Order("created_at DESC").Find(&users).Error; err != nil {
		return users, err
	}
	return users, nil
}

func IncreaseUserMessagesStats(user *models.User) {
	user.LastMessageDate = time.Now()
	user.MessageCount += 1
	database.DB.Save(user)
}
