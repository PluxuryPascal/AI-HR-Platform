package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseResult_ParsingStatus(t *testing.T) {
	tests := []struct {
		name    string
		missing []string
		want    CandidateParsingStatus
	}{
		{name: "все поля на месте", missing: nil, want: ParsingStatusCompleted},
		{name: "не хватает необязательного поля", missing: []string{"education", "summary"}, want: ParsingStatusCompleted},
		{name: "нет email", missing: []string{"email"}, want: ParsingStatusNeedsReview},
		{name: "нет имени", missing: []string{"first_name"}, want: ParsingStatusNeedsReview},
		{name: "нет фамилии", missing: []string{"last_name"}, want: ParsingStatusNeedsReview},
		{name: "нет локации среди прочих", missing: []string{"summary", "location"}, want: ParsingStatusNeedsReview},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ParseResult{MissingFields: tt.missing}
			assert.Equal(t, tt.want, r.ParsingStatus())
		})
	}
}
