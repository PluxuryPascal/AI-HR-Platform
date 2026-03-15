package pdf

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/ledongthuc/pdf"
)

var (
	// ErrPasswordProtected — PDF зашифрован паролем. Повтор бесполезен.
	ErrPasswordProtected = errors.New("pdf: password protected")

	// ErrImageBasedPDF — PDF содержит только отсканированные изображения без
	// текстового слоя. Требует LLM Vision / OCR fallback.
	ErrImageBasedPDF = errors.New("pdf: image-based, no text layer found")

	// ErrEmptyContent — файл пуст или повреждён настолько, что текст не читается.
	ErrEmptyContent = errors.New("pdf: empty or corrupted content")
)

type Extractor interface {
	Extract(data []byte) (string, error)
}

// Config настраивает поведение экстрактора.
type Config struct {
	// MinTextLength — минимальное количество значимых символов после извлечения.
	// Если меньше — PDF считается image-based и возвращается ErrImageBasedPDF.
	// По умолчанию: 100.
	MinTextLength int

	// MaxPages — максимальное количество страниц для обработки (0 = все).
	// Защита от аномально больших файлов.
	MaxPages int
}

func (c *Config) applyDefaults() {
	if c.MinTextLength == 0 {
		c.MinTextLength = 100
	}
}

type extractor struct {
	cfg Config
}

// NewExtractor создаёт экстрактор с заданной конфигурацией.
func NewExtractor(cfg Config) Extractor {
	cfg.applyDefaults()
	return &extractor{cfg: cfg}
}

// Extract извлекает текст из PDF-байтов.
//
// Шаги:
//  1. Обнаружение пароля.
//  2. Парсинг через ledongthuc/pdf — читает текстовый слой постранично.
//  3. Сборка и нормализация текста.
//  4. Проверка минимальной длины → ErrImageBasedPDF если короче порога.
func (e *extractor) Extract(data []byte) (string, error) {
	if len(data) == 0 {
		return "", ErrEmptyContent
	}

	if isPasswordProtected(data) {
		return "", ErrPasswordProtected
	}

	r := bytes.NewReader(data)

	pdfReader, err := pdf.NewReader(r, r.Size())
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrEmptyContent, err)
	}

	totalPages := pdfReader.NumPage()
	if totalPages == 0 {
		return "", ErrEmptyContent
	}

	limit := totalPages
	if e.cfg.MaxPages > 0 && e.cfg.MaxPages < totalPages {
		limit = e.cfg.MaxPages
	}

	var sb strings.Builder
	sb.Grow(len(data) / 4)

	for pageNum := 1; pageNum <= limit; pageNum++ {
		page := pdfReader.Page(pageNum)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			continue
		}

		if text != "" {
			sb.WriteString(text)
			sb.WriteByte('\n')
		}
	}

	raw := sb.String()
	cleaned := cleanText(raw)

	if countMeaningfulChars(cleaned) < e.cfg.MinTextLength {
		// Текст слишком короткий — PDF либо image-based, либо содержит только
		// мусор (артефакты шрифтов, управляющие символы).
		return "", ErrImageBasedPDF
	}

	return cleaned, nil
}

// isPasswordProtected проверяет наличие /Encrypt словаря в PDF-заголовке.
func isPasswordProtected(data []byte) bool {
	header := data
	if len(header) > 8192 {
		header = header[:8192]
	}
	return bytes.Contains(header, []byte("/Encrypt"))
}

// cleanText нормализует извлечённый текст:
//   - убирает нулевые байты и управляющие символы (кроме \n, \t)
//   - схлопывает множественные пробелы в один
//   - схлопывает больше двух переводов строк подряд в два
//   - обрезает пробелы по краям
func cleanText(s string) string {
	if s == "" {
		return ""
	}

	var sb strings.Builder
	sb.Grow(len(s))

	prevSpace := false
	newlineCount := 0

	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		i += size

		if r == utf8.RuneError && size == 1 {
			prevSpace = false
			continue
		}

		switch {
		case r == '\n':
			newlineCount++
			prevSpace = false
			if newlineCount <= 2 {
				sb.WriteRune('\n')
			}

		case r == '\r':

		case r == '\t':
			if !prevSpace {
				sb.WriteByte(' ')
				prevSpace = true
			}
			newlineCount = 0

		case unicode.IsControl(r):

		case unicode.IsSpace(r):
			newlineCount = 0
			if !prevSpace {
				sb.WriteByte(' ')
				prevSpace = true
			}

		default:
			newlineCount = 0
			prevSpace = false
			sb.WriteRune(r)
		}
	}

	return strings.TrimSpace(sb.String())
}

// countMeaningfulChars считает буквы и цифры в строке.
// Используется для отличия реального текста от мусора (артефактов, пробелов).
func countMeaningfulChars(s string) int {
	count := 0
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			count++
		}
	}
	return count
}
