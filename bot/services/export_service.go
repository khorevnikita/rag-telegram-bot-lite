package services

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"gorag-telegram-bot/models"
)

const DEFAULT_SHEET_NAME = "Sheet1"

func addHeaders(f *excelize.File, headers []string) error {
	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+i)))
		err := f.SetCellValue(DEFAULT_SHEET_NAME, cell, header)
		if err != nil {
			return err
		}
	}
	return nil
}

func ExportUsersToExcel(users []models.User, questions []models.Question, answers []models.Answer) (*bytes.Buffer, error) {
	// Создаем новый Excel файл
	f := excelize.NewFile()

	// Создаем базовые заголовки
	baseHeaders := []string{"Telegram ID", "Telegram Username", "First Name", "Last Name", "Created At", "Messages Count", "Last Message Date", "Form Completed At"}
	_ = addHeaders(f, baseHeaders)

	// Добавляем вопросы как заголовки
	for i, question := range questions {
		cell := fmt.Sprintf("%s1", string(rune('A'+len(baseHeaders)+i)))
		_ = f.SetCellValue(DEFAULT_SHEET_NAME, cell, question.Text)
	}

	// Создаем маппинг ответов для быстрого поиска
	answerMap := make(map[uuid.UUID]map[uuid.UUID]string) // map[UserID][QuestionID]AnswerText
	for _, answer := range answers {
		if _, ok := answerMap[*answer.UserId]; !ok {
			answerMap[*answer.UserId] = make(map[uuid.UUID]string)
		}
		answerMap[*answer.UserId][answer.QuestionId] = answer.Text
	}

	// Добавляем данные пользователей и ответы
	for i, user := range users {
		row := i + 2 // Первая строка — это заголовки
		_ = f.SetCellValue(DEFAULT_SHEET_NAME, fmt.Sprintf("A%d", row), user.TelegramID)
		_ = f.SetCellValue(DEFAULT_SHEET_NAME, fmt.Sprintf("B%d", row), dereferenceString(user.TelegramUsername))
		_ = f.SetCellValue(DEFAULT_SHEET_NAME, fmt.Sprintf("C%d", row), dereferenceString(user.FirstName))
		_ = f.SetCellValue(DEFAULT_SHEET_NAME, fmt.Sprintf("D%d", row), dereferenceString(user.LastName))
		_ = f.SetCellValue(DEFAULT_SHEET_NAME, fmt.Sprintf("E%d", row), user.CreatedAt.Format("2006-01-02 15:04:05"))
		_ = f.SetCellValue(DEFAULT_SHEET_NAME, fmt.Sprintf("F%d", row), user.MessageCount)
		_ = f.SetCellValue(DEFAULT_SHEET_NAME, fmt.Sprintf("G%d", row), user.LastMessageDate.Format("2006-01-02 15:04:05"))
		if user.FormCompletedAt != nil {
			_ = f.SetCellValue(DEFAULT_SHEET_NAME, fmt.Sprintf("H%d", row), user.FormCompletedAt.Format("2006-01-02 15:04:05"))
		}

		// Добавляем ответы пользователя
		for j, question := range questions {
			col := string(rune('A' + len(baseHeaders) + j))
			answerText := ""
			if userAnswers, ok := answerMap[user.ID]; ok {
				if text, exists := userAnswers[question.ID]; exists {
					answerText = text
				}
			}
			_ = f.SetCellValue(DEFAULT_SHEET_NAME, fmt.Sprintf("%s%d", col, row), answerText)
		}
	}

	// Сохраняем файл в буфер
	buf := new(bytes.Buffer)
	if err := f.Write(buf); err != nil {
		return nil, fmt.Errorf("failed to write file to buffer: %w", err)
	}

	return buf, nil
}

func ExportDislikedMessages(messages []models.Message) (*bytes.Buffer, error) {
	// Создаем новый Excel файл
	f := excelize.NewFile()

	// Создаем заголовки
	headers := []string{
		"Message ID", "Username", "Content", "Response", "Created At", "Updated At",
	}
	err := addHeaders(f, headers)
	if err != nil {
		return nil, err
	}

	// Заполняем данные
	for i, message := range messages {
		row := i + 2 // Первая строка — это заголовки
		_ = f.SetCellValue(DEFAULT_SHEET_NAME, fmt.Sprintf("A%d", row), message.ID)
		_ = f.SetCellValue(DEFAULT_SHEET_NAME, fmt.Sprintf("B%d", row), dereferenceString(message.User.TelegramUsername))
		_ = f.SetCellValue(DEFAULT_SHEET_NAME, fmt.Sprintf("C%d", row), message.Content)
		_ = f.SetCellValue(DEFAULT_SHEET_NAME, fmt.Sprintf("D%d", row), dereferenceString(message.Response))
		_ = f.SetCellValue(DEFAULT_SHEET_NAME, fmt.Sprintf("E%d", row), message.CreatedAt.Format("2006-01-02 15:04:05"))
		_ = f.SetCellValue(DEFAULT_SHEET_NAME, fmt.Sprintf("F%d", row), message.UpdatedAt.Format("2006-01-02 15:04:05"))
	}

	// Сохраняем файл в буфер
	buf := new(bytes.Buffer)
	if err := f.Write(buf); err != nil {
		return nil, fmt.Errorf("failed to write file to buffer: %w", err)
	}

	return buf, nil
}
func ExportSubscriptions(data []SubscriptionReport) (*bytes.Buffer, error) {
	// Создаем новый Excel файл
	f := excelize.NewFile()

	// Создаем заголовки
	headers := []string{
		"User ID", "Telegram ID", "Username", "Subscribed At", "Expires At", "Unsubscribed At", "Money spent", "Last payment date",
	}
	err := addHeaders(f, headers)
	if err != nil {
		return nil, err
	}

	// Заполняем данные
	for i, item := range data {
		row := i + 2 // Первая строка — это заголовки
		_ = f.SetCellValue(DEFAULT_SHEET_NAME, fmt.Sprintf("A%d", row), item.UserID)
		_ = f.SetCellValue(DEFAULT_SHEET_NAME, fmt.Sprintf("B%d", row), item.TelegramID)
		_ = f.SetCellValue(DEFAULT_SHEET_NAME, fmt.Sprintf("C%d", row), dereferenceString(item.TelegramUsername))
		_ = f.SetCellValue(DEFAULT_SHEET_NAME, fmt.Sprintf("D%d", row), item.SubscriptionCreatedAt.Format("2006-01-02 15:04:05"))
		_ = f.SetCellValue(DEFAULT_SHEET_NAME, fmt.Sprintf("E%d", row), item.SubscriptionExpiresAt.Format("2006-01-02 15:04:05"))
		if item.UnsubscribedAt != nil {
			_ = f.SetCellValue(DEFAULT_SHEET_NAME, fmt.Sprintf("F%d", row), item.UnsubscribedAt.Format("2006-01-02 15:04:05"))
		}
		_ = f.SetCellValue(DEFAULT_SHEET_NAME, fmt.Sprintf("G%d", row), item.PaymentsTotalAmount)
		if item.LastPaymentAt != nil {
			_ = f.SetCellValue(DEFAULT_SHEET_NAME, fmt.Sprintf("H%d", row), item.LastPaymentAt.Format("2006-01-02 15:04:05"))
		}
	}

	// Сохраняем файл в буфер
	buf := new(bytes.Buffer)
	if err := f.Write(buf); err != nil {
		return nil, fmt.Errorf("failed to write file to buffer: %w", err)
	}

	return buf, nil
}

// Вспомогательная функция для разыменования указателя на строку
func dereferenceString(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}
