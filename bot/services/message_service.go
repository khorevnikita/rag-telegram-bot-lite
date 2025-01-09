package services

import (
	"fmt"
	tb "gopkg.in/telebot.v4"
	"gorag-telegram-bot/database"
	"gorag-telegram-bot/models"
)

func SaveMessage(user *models.User, tempMsg *tb.Message, text string, files []QuestionFile) (models.Message, error) {
	// Сохраняем сообщение пользователя в базе данных
	messageLog := models.Message{
		UserID:            user.ID,
		TelegramMessageID: tempMsg.ID,
		Content:           text,
		SystemContextID:   user.SystemContextID,
	}
	if err := database.DB.Create(&messageLog).Error; err != nil {
		fmt.Printf("Error saving message log: %v\n", err)
		return messageLog, err
	}

	var messageFiles []models.MessageFile
	for _, f := range files {
		messageFiles = append(messageFiles, models.MessageFile{
			MessageID: messageLog.ID,
			UserID:    user.ID,
			FilePath:  f.FilePath,
			FileName:  f.Filename,
			Size:      f.Size,
			Extension: f.Extension,
			FileType:  f.FileType,
		})
	}

	if len(messageFiles) > 0 {
		err := database.DB.Create(&messageFiles).Error
		if err != nil {
			return messageLog, err
		}
	}

	return messageLog, nil
}

func UpdateMessage(messageLog *models.Message, aiQuestion *AIQuestion) {
	// Обновляем поле ответа в записи messageLog
	messageLog.Response = &aiQuestion.Answer
	messageLog.AIMessageID = &aiQuestion.ID // Сохраняем AIMessageID
	database.DB.Save(&messageLog)
}

func FindByTelegramID(tgID int) (*models.Message, error) {
	var message models.Message
	if err := database.DB.Where("telegram_message_id = ?", tgID).First(&message).Error; err != nil {
		return nil, err
	}
	return &message, nil
}

func SetLiked(message *models.Message) error {
	// Обновляем оценку в базе данных
	like := true
	if err := database.DB.Model(message).Update("liked", &like).Error; err != nil {
		return err
	}
	return nil
}
func SetDisliked(message *models.Message) error {
	// Обновляем оценку в базе данных
	dislike := false
	if err := database.DB.Model(message).Update("liked", &dislike).Error; err != nil {
		return err
	}
	return nil
}

// GetDislikedMessages возвращает список сообщений с отрицательной оценкой
func GetDislikedMessages() ([]models.Message, error) {
	var dislikedMessages []models.Message
	err := database.DB.
		Preload("User").
		Where("liked = ?", false).
		Order("created_at desc").
		Find(&dislikedMessages).Error
	if err != nil {
		return nil, err
	}
	return dislikedMessages, nil
}

func ChangeTelegramId(messageLog *models.Message, tgID int) {
	// Обновляем поле ответа в записи messageLog
	messageLog.TelegramMessageID = tgID
	database.DB.Save(&messageLog)
}
