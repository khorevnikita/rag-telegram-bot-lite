package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gorag-telegram-bot/config"
	"gorag-telegram-bot/models"
	"io"
	"mime/multipart"
	"net/http"
)

type QuestionFile struct {
	FilePath  string `json:"file_path"`
	Filename  string `json:"filename"`
	Size      int    `json:"size"`
	Extension string `json:"extension"`
	FileType  string `json:"file_type"`
}

type AIQuestion struct {
	Answer         string `json:"answer"`
	ConversationID int    `json:"conversation_id"`
	ID             int    `json:"id"`
}

type AIResponse struct {
	Question AIQuestion `json:"question"`
}

type AIClient struct {
	URL      string
	APIToken string
}

func NewAIClient() AIClient {
	appConf := config.AppConfig
	return AIClient{
		URL:      appConf.RagApiEndpoint,
		APIToken: appConf.RagApiToken,
	}
}

func (c AIClient) sendAIRequest(requestType string, path string, body map[string]interface{}) ([]byte, error) {
	bodyJSON, _ := json.Marshal(body)
	fmt.Printf("DEBUG: %s // %s // %s\n", c.URL, c.APIToken, bodyJSON)
	req, _ := http.NewRequest(requestType, fmt.Sprintf("%s/api%s", c.URL, path), bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.APIToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	// Читаем тело ответа в виде строки
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	fmt.Printf("DEBUG: StatusCode=%d, ResponseBody=%s\n", resp.StatusCode, string(bodyBytes))
	return bodyBytes, nil
}

func (c AIClient) Answer(conversationID *int, request string, files []QuestionFile, context string) (*AIQuestion, error) {
	dialogId := 0
	if conversationID != nil {
		dialogId = *conversationID
	}

	body := map[string]interface{}{
		"text":            request,
		"conversation_id": dialogId,
		"stream":          false,
		"type":            "text",
		"answer_type":     "text",
		"webhook":         "",
		"context":         context,
		"files":           files,
	}

	fmt.Printf("BODY REQ: %+v\n", body)

	response, err := c.sendAIRequest("POST", "/questions", body)

	if err != nil {
		return nil, err
	}

	// Декодируем JSON-ответ
	var aiResp AIResponse
	if err := json.Unmarshal(response, &aiResp); err != nil {
		return nil, err
	}

	return &aiResp.Question, nil
}

func (c AIClient) LikeQuestion(questionID int) error {
	_, err := c.sendAIRequest("POST", fmt.Sprintf("/questions/%d/like", questionID), map[string]interface{}{})
	return err
}

func (c AIClient) DislikeQuestion(questionID int) error {
	_, err := c.sendAIRequest("POST", fmt.Sprintf("/questions/%d/dislike", questionID), map[string]interface{}{})
	return err
}

func (c AIClient) Upload(fileData []byte, fileName string) (*QuestionFile, error) {
	url := fmt.Sprintf("%s/api/storage/upload", c.URL)

	// Создаём тело form-data запроса
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Добавляем файл в form-data
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать form-data для файла: %w", err)
	}
	if _, err = part.Write(fileData); err != nil {
		return nil, fmt.Errorf("не удалось записать файл в form-data: %w", err)
	}

	// Завершаем формирование тела запроса
	writer.Close()

	// Создаём HTTP запрос
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать HTTP запрос: %w", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("x-api-key", c.APIToken)

	// Отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("не удалось загрузить файл, статус: %d, тело: %s", resp.StatusCode, string(responseBody))
	}

	// Читаем тело ответа
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать ответ: %w", err)
	}

	// Возвращаем тело ответа (например, ID загруженного файла)
	fmt.Printf("DEBUG: Успешный ответ: %s\n", responseBody)
	var uploadedFile QuestionFile
	if err := json.Unmarshal(responseBody, &uploadedFile); err != nil {
		return nil, err
	}
	return &uploadedFile, nil
}

func (c AIClient) GetAIResponse(user *models.User, messageLog *models.Message, files []QuestionFile) (*AIQuestion, error) {
	var context = ""

	if config.AppConfig.Modules.Form.Enabled {
		context += fmt.Sprintf("%s\n%s\n", config.AppConfig.Modules.Form.ContextPrefix, SerializeUserAnswers(user))
	}

	if messageLog.SystemContextID != nil {
		sysCtx, _ := FindSystemContext(messageLog.SystemContextID)
		if sysCtx != nil {
			context += fmt.Sprintf("%s\n", sysCtx.Text)
		}
	}

	aiQuestion, err := c.Answer(user.ConversationID, messageLog.Content, files, context)
	if err != nil || aiQuestion == nil {
		return nil, err
	}

	// Обновляем ConversationID, если это первый запрос
	if user.ConversationID == nil {
		SetUserConversation(user, &aiQuestion.ConversationID)
	}

	UpdateMessage(messageLog, aiQuestion)
	return aiQuestion, nil
}
