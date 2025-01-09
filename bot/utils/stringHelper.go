package utils

import (
	"golang.org/x/exp/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func PointerToString(s string) *string {
	return &s
}

// Функция для разбиения текста на части
func SplitTextIntoChunks(text string, maxLen int) []string {
	var chunks []string

	for len(text) > maxLen {
		// Ищем ближайший перенос строки в пределах maxLen
		cut := strings.LastIndex(text[:maxLen], "\n")
		if cut == -1 {
			// Если новой строки нет, ищем пробел
			cut = strings.LastIndex(text[:maxLen], " ")
		}
		if cut == -1 || cut == 0 { // Если не нашли перенос строки или пробел, делим строго по длине
			cut = maxLen
		}

		chunks = append(chunks, text[:cut])
		text = text[cut:]
	}

	chunks = append(chunks, text)
	return chunks
}

// Вспомогательная функция для преобразования строки в float64
func ParseFloat(value string) float64 {
	f, _ := strconv.ParseFloat(value, 64)
	return f
}

func ParseInt(value string) int {
	f, err := strconv.ParseInt(value, 10, 0) // Основание 10, размер 0 (определяется автоматически)
	if err != nil {
		return 0
	}
	return int(f)
}

// Вспомогательная функция для декодирования URL-кодированных строк
func DecodeURLEncoding(value string) string {
	decoded, _ := url.QueryUnescape(value)
	return decoded
}

// generateRandomPassword генерирует случайный пароль
func GenerateRandomPassword() string {
	const passwordLength = 12
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()"
	rand.Seed(uint64(time.Now().UnixNano()))
	password := make([]byte, passwordLength)
	for i := range password {
		password[i] = charset[rand.Intn(len(charset))]
	}
	return string(password)
}
func Coalesce(sp *string, s string) string {
	if sp == nil {
		return s
	}
	return *sp
}
