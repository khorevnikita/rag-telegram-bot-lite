package handlers

import (
	"bytes"
	"fmt"
	tb "gopkg.in/telebot.v4"
	"gorag-telegram-bot/config"
	"gorag-telegram-bot/core"
	"gorag-telegram-bot/models"
	"gorag-telegram-bot/services"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type InputHandlers struct {
	AIClient    services.AIClient
	formHandler FormHandlers
}

func GetNewInputHandlers(formHandler FormHandlers) InputHandlers {

	return InputHandlers{
		AIClient:    services.NewAIClient(),
		formHandler: formHandler,
	}
}

func (h InputHandlers) RegisterCommands() {
	core.Bot.Handle(tb.OnText, h.onInput)

	moduleConf := config.AppConfig.Modules

	if moduleConf.Files.Enabled {
		core.Bot.Handle(tb.OnDocument, h.onInput)
	}

	if moduleConf.Audio.Enabled {
		core.Bot.Handle(tb.OnAudio, h.onInput)
		core.Bot.Handle(tb.OnMedia, h.onInput)
		core.Bot.Handle(tb.OnVoice, h.onInput)
	}
}

func (h InputHandlers) onInput(c tb.Context) error {
	user, err := services.GetUserByTelegramID(c.Sender().ID)
	if err != nil {
		return h.responseNoUser(c)
	}

	if user.State != nil {
		err = h.processState(c, user)
		if err != nil {
			return c.Send(err.Error())
		}
		return nil
	}

	modulesConfig := config.AppConfig.Modules
	if modulesConfig.Form.Enabled && !modulesConfig.Form.CanSkip && user.FormCompletedAt == nil {
		return services.SendStartFormMessage(c.Sender().ID)
	}

	if modulesConfig.Billing.Enabled {
		subscription, err := services.GetSubscription(user)
		if err != nil {
			return c.Send(err.Error())
		}
		if subscription == nil || !subscription.IsActive() {
			return services.SendSubscribeMessage(c)
		}
	}

	services.IncreaseUserMessagesStats(user)

	tempMsg, err := services.SendTemporaryMessage(c)
	if err != nil {
		return err
	}

	allFiles, err := h.collectFilesFromContext(c)
	if err != nil {
		return h.processErrorMessage(c, tempMsg)
	}

	messageLog, err := services.SaveMessage(user, tempMsg, c.Text(), allFiles)
	if err != nil {
		return h.processErrorMessage(c, tempMsg)
	}

	_, err = h.AIClient.GetAIResponse(user, &messageLog, allFiles)

	if err != nil {
		return h.processErrorMessage(c, tempMsg)
	}

	return services.SendResponseMessage(c, tempMsg, &messageLog)
}

func (h InputHandlers) collectFilesFromContext(c tb.Context) ([]services.QuestionFile, error) {
	var allFiles []services.QuestionFile

	file, err := h.processFiles(c)
	if err != nil {
		return nil, err
	}
	if file != nil {
		allFiles = append(allFiles, *file)
	}

	audio, err := h.processAudio(c)
	if err != nil {
		return nil, err
	}
	if audio != nil {
		allFiles = append(allFiles, *audio)
	}

	voice, err := h.processVoice(c)
	if err != nil {
		return nil, err
	}

	if voice != nil {
		allFiles = append(allFiles, *voice)
	}

	video, err := h.processVideo(c)
	if err != nil {
		return nil, err
	}
	if video != nil {
		allFiles = append(allFiles, *video)
	}
	return allFiles, nil
}

func (h InputHandlers) responseNoUser(c tb.Context) error {
	return c.Send("Пользователь не найден. Пожалуйста, начните с команды /start.")
}

func (h InputHandlers) processErrorMessage(c tb.Context, tempMsg *tb.Message) error {
	// Если произошла ошибка, обновляем сообщение с текстом ошибки
	errorMessage := "Что-то пошло не так. Попробуйте еще раз. Если ошибка повторяется, свяжитесь с тех. поддержкой"
	_, err := c.Bot().Edit(tempMsg, errorMessage)
	return err
}

func (h InputHandlers) transferFile(fileName string, file tb.File) (*services.QuestionFile, error) {
	data, err := h.downloadTelegramFile(file.FileID)
	if err != nil {
		fmt.Printf("Ошибка скачивании фото: %v\n", err)
		return nil, err
	}
	fmt.Printf("DEBUG FILE %s (%d): %s\n", file.FileURL, file.FileSize, fileName)

	uploadedFile, err := h.AIClient.Upload(data, fileName)
	if err != nil {
		fmt.Printf("Ошибка загрузки фото: %v\n", err)
		return nil, err
	}

	fmt.Printf("Uploaded %+v , %+v\n", uploadedFile, err)
	return uploadedFile, nil
}

func (h InputHandlers) processFiles(c tb.Context) (*services.QuestionFile, error) {
	document := c.Message().Document
	if document != nil {
		return h.transferFile(document.FileName, document.File)
	}
	return nil, nil
}

func (h InputHandlers) processAudio(c tb.Context) (*services.QuestionFile, error) {
	audio := c.Message().Audio
	if audio != nil {
		return h.transferFile(audio.FileName, audio.File)
	}
	return nil, nil
}

func (h InputHandlers) processVideo(c tb.Context) (*services.QuestionFile, error) {
	audio := c.Message().Video
	if audio != nil {
		return h.transferFile(audio.FileName, audio.File)
	}
	return nil, nil
}

func (h InputHandlers) processVoice(c tb.Context) (*services.QuestionFile, error) {
	voice := c.Message().Voice
	if voice != nil {
		fmt.Printf("DEBUG VOICE: mime: %s, caption: %s, url: %s, path: %s\n", voice.MIME, voice.Caption, voice.File.FileURL, voice.File.FilePath)
		// Генерируем уникальные имена файлов
		timestamp := time.Now().UnixNano()
		inputFileName := fmt.Sprintf("input_%d.oga", timestamp)
		outputFileName := fmt.Sprintf("output_%d.wav", timestamp) // Или .wav

		// Скачиваем голосовой файл
		data, err := h.downloadTelegramFile(voice.FileID)
		if err != nil {
			fmt.Printf("error downloading %s\n", err)
			return nil, err
		}

		err = os.WriteFile(inputFileName, data, 0644)
		if err != nil {
			fmt.Printf("error writing %s\n", err)
			return nil, err
		}

		// Конвертируем файл с помощью ffmpeg
		cmd := exec.Command("ffmpeg", "-i", inputFileName, "-ar", "16000", "-ac", "1", "-sample_fmt", "s16", outputFileName)
		//cmd := exec.Command("ffmpeg", "-i", inputFileName, "-ar", "16000", "-ac", "1", outputFileName)
		//cmd := exec.Command("ffmpeg", "-i", inputFileName, "-acodec", "copy", outputFileName)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		err = cmd.Run()
		if err != nil {
			fmt.Printf("error convertation %s, %s\n", err, stderr.String())
			return nil, fmt.Errorf("Ошибка при конвертации ffmpeg: %v: %s", err, stderr.String())
		}
		defer os.Remove(outputFileName) // Удаляем выходной файл после обработки

		convertedFile, err := os.Open(outputFileName)
		if err != nil {
			fmt.Printf("error opening %s\n", err)
			return nil, err
		}
		defer convertedFile.Close()

		// Читаем содержимое файла в байтовый срез
		convertedFileBytes, err := io.ReadAll(convertedFile)
		if err != nil {
			fmt.Printf("error byting %s\n", err)
			return nil, err
		}

		return h.AIClient.Upload(convertedFileBytes, outputFileName)
	}
	return nil, nil
}

func (h InputHandlers) downloadTelegramFile(fileID string) ([]byte, error) {
	// Получаем объект файла, чтобы получить FilePath
	file, err := core.Bot.FileByID(fileID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить путь файла: %w", err)
	}
	fmt.Printf("DEBUG: FilePath: %s\n", file.FilePath)

	// Формируем URL для загрузки файла
	baseURL := "https://api.telegram.org/file/bot"
	url := fmt.Sprintf("%s%s/%s", baseURL, core.Bot.Token, file.FilePath)

	// Отправляем GET-запрос
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("не удалось загрузить файл: %w", err)
	}
	defer resp.Body.Close()
	fmt.Printf("DEBUG URL: %s\n", url)
	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("не удалось загрузить файл, статус: %d, тело: %s", resp.StatusCode, string(body))
	}

	// Читаем содержимое файла
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать содержимое файла: %w", err)
	}

	fmt.Printf("DEBUG: Размер загруженного файла: %d байт\n", len(data))
	return data, nil
}

func (h InputHandlers) processState(c tb.Context, user *models.User) error {
	switch *user.State {
	case "form":
		return h.formHandler.ProcessFormAnswer(c, user)
	}
	// never should be here
	fmt.Printf("NEVER SHOULD BEEN HERE ERROR\n")
	return fmt.Errorf("внутренняя ошибка. Вопрос не найден. Обратитесь в тех. поддержку")
}
