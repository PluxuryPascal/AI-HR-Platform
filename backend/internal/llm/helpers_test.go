package llm

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSplitIntoChunks(t *testing.T) {
	t.Run("пустой текст", func(t *testing.T) {
		assert.Nil(t, splitIntoChunks("", 10, 2))
	})

	t.Run("текст короче чанка возвращается целиком", func(t *testing.T) {
		assert.Equal(t, []string{"hello world"}, splitIntoChunks("hello world", 100, 10))
	})

	t.Run("длинный текст режется по границе слова", func(t *testing.T) {
		text := strings.Repeat("word ", 200) // 1000 байт
		chunks := splitIntoChunks(text, 100, 20)

		require.Greater(t, len(chunks), 1)
		for i, c := range chunks {
			assert.LessOrEqual(t, len(c), 100, "чанк %d длиннее лимита", i)
			if i < len(chunks)-1 {
				assert.False(t, strings.HasSuffix(c, "wor"), "чанк %d обрезан посреди слова", i)
			}
		}
	})

	t.Run("соседние чанки перекрываются", func(t *testing.T) {
		text := strings.Repeat("abcd ", 100)
		chunks := splitIntoChunks(text, 50, 15)

		require.Greater(t, len(chunks), 1)
		tail := chunks[0][len(chunks[0])-10:]
		assert.Contains(t, chunks[1], strings.TrimSpace(tail))
	})

	t.Run("весь текст покрыт чанками", func(t *testing.T) {
		text := strings.Repeat("lorem ipsum dolor ", 60)
		chunks := splitIntoChunks(text, 128, 32)

		assert.True(t, strings.HasPrefix(text, chunks[0]))
		assert.True(t, strings.HasSuffix(text, chunks[len(chunks)-1]))
	})

	t.Run("нулевой overlap не зацикливается", func(t *testing.T) {
		text := strings.Repeat("x", 300) // без пробелов
		chunks := splitIntoChunks(text, 100, 0)

		assert.Equal(t, 3, len(chunks))
		assert.Equal(t, text, strings.Join(chunks, ""))
	})
}

func TestFindWordBoundary(t *testing.T) {
	tests := []struct {
		name string
		text string
		pos  int
		want int
	}{
		{name: "позиция за концом текста", text: "abc", pos: 10, want: 3},
		{name: "пробел ровно на позиции", text: "aaa bbb", pos: 3, want: 3},
		{name: "ищет пробел левее позиции", text: "aaa bbbbbb", pos: 7, want: 3},
		{name: "перенос строки тоже граница", text: "aaa\nbbbbbb", pos: 7, want: 3},
		{name: "без пробелов возвращает позицию", text: strings.Repeat("a", 50), pos: 20, want: 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, findWordBoundary(tt.text, tt.pos))
		})
	}
}

func TestSanitise(t *testing.T) {
	r := &JobParseResult{
		Title:        "  Go Developer \n",
		Description:  "\tОписание  ",
		Currency:     " rub ",
		WorkFormat:   "remote", // невалидный регистр
		Requirements: []string{" Go ", "", "   ", "PostgreSQL"},
		SalaryMin:    -100,
		SalaryMax:    -1,
	}

	sanitise(r)

	assert.Equal(t, "Go Developer", r.Title)
	assert.Equal(t, "Описание", r.Description)
	assert.Equal(t, "RUB", r.Currency)
	assert.Equal(t, "Onsite", r.WorkFormat)
	assert.Equal(t, []string{"Go", "PostgreSQL"}, r.Requirements)
	assert.Zero(t, r.SalaryMin)
	assert.Zero(t, r.SalaryMax)
}

func TestSanitise_KeepsValidWorkFormat(t *testing.T) {
	for _, f := range []string{"Remote", "Onsite", "Hybrid"} {
		r := &JobParseResult{WorkFormat: f}
		sanitise(r)
		assert.Equal(t, f, r.WorkFormat)
	}
}

func TestLocaleToLanguage(t *testing.T) {
	tests := map[string]string{
		"ru":    "Russian",
		"ru-RU": "Russian",
		"es":    "Spanish",
		"es-MX": "Spanish",
		"en":    "English",
		"de":    "English",
		"":      "English",
	}

	for locale, want := range tests {
		t.Run(locale, func(t *testing.T) {
			assert.Equal(t, want, localeToLanguage(locale))
		})
	}
}
