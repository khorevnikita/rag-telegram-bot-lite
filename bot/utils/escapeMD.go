package utils

import (
	"regexp"
	"strings"
)

// Экранирование текста для MarkdownV2 с обработкой заголовков
func EscapeMarkdownV2WithHeaders(text string) string {
	// Регулярное выражение для заголовков
	headerRegex := regexp.MustCompile(`^(#+)\s*(.*)`)

	// Обработка каждой строки текста
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if headerRegex.MatchString(line) {
			matches := headerRegex.FindStringSubmatch(line)
			headerContent := matches[2] // Содержимое заголовка после `#`
			// Экранируем только содержимое заголовка
			lines[i] = "*" + escapeMarkdownV2Content(headerContent) + "*"
		} else {
			// Экранируем остальные строки, не затрагивая существующее форматирование
			lines[i] = escapeMarkdownV2Content(line)
		}
	}

	// Собираем обработанный текст обратно
	return strings.Join(lines, "\n")
}

// Экранирование текста для MarkdownV2 с сохранением существующего форматирования
func escapeMarkdownV2Content(text string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		".", "\\.",
		"!", "\\!",
		"`", "\\`",
		"\\", "\\\\",
	)
	return replacer.Replace(text)
}
