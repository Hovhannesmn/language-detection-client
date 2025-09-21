package validator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"language-detection-client/internal/models"
)

func TestCoverageValidator_Validate(t *testing.T) {
	validator := NewCoverageValidator()

	tests := []struct {
		name           string
		captions       []models.Caption
		config         models.Config
		expectedErrors int
		expectedTypes  []string
	}{
		{
			name: "sufficient coverage",
			captions: []models.Caption{
				{
					Start: 1 * time.Second,
					End:   5 * time.Second,
					Text:  "Hello world",
				},
				{
					Start: 6 * time.Second,
					End:   10 * time.Second,
					Text:  "Goodbye world",
				},
			},
			config: models.Config{
				StartTime:        0,
				EndTime:          10 * time.Second,
				CoverageRequired: 80.0,
			},
			expectedErrors: 0,
			expectedTypes:  []string{},
		},
		{
			name: "insufficient coverage",
			captions: []models.Caption{
				{
					Start: 1 * time.Second,
					End:   3 * time.Second,
					Text:  "Hello world",
				},
			},
			config: models.Config{
				StartTime:        0,
				EndTime:          10 * time.Second,
				CoverageRequired: 80.0,
			},
			expectedErrors: 1,
			expectedTypes:  []string{"caption_coverage"},
		},
		{
			name: "invalid time range - end before start",
			captions: []models.Caption{
				{
					Start: 1 * time.Second,
					End:   3 * time.Second,
					Text:  "Hello world",
				},
			},
			config: models.Config{
				StartTime:        10 * time.Second,
				EndTime:          5 * time.Second,
				CoverageRequired: 80.0,
			},
			expectedErrors: 1,
			expectedTypes:  []string{"invalid_time_range"},
		},
		{
			name: "invalid time range - same start and end",
			captions: []models.Caption{
				{
					Start: 1 * time.Second,
					End:   3 * time.Second,
					Text:  "Hello world",
				},
			},
			config: models.Config{
				StartTime:        5 * time.Second,
				EndTime:          5 * time.Second,
				CoverageRequired: 80.0,
			},
			expectedErrors: 1,
			expectedTypes:  []string{"invalid_time_range"},
		},
		{
			name:     "no captions",
			captions: []models.Caption{},
			config: models.Config{
				StartTime:        0,
				EndTime:          10 * time.Second,
				CoverageRequired: 80.0,
			},
			expectedErrors: 1,
			expectedTypes:  []string{"caption_coverage"},
		},
		{
			name: "captions outside time range",
			captions: []models.Caption{
				{
					Start: 15 * time.Second,
					End:   20 * time.Second,
					Text:  "Hello world",
				},
			},
			config: models.Config{
				StartTime:        0,
				EndTime:          10 * time.Second,
				CoverageRequired: 80.0,
			},
			expectedErrors: 1,
			expectedTypes:  []string{"caption_coverage"},
		},
		{
			name: "partial overlap with time range",
			captions: []models.Caption{
				{
					Start: 5 * time.Second,
					End:   15 * time.Second,
					Text:  "Hello world",
				},
			},
			config: models.Config{
				StartTime:        0,
				EndTime:          10 * time.Second,
				CoverageRequired: 80.0,
			},
			expectedErrors: 1,
			expectedTypes:  []string{"caption_coverage"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validator.Validate(tt.captions, tt.config)
			assert.NoError(t, err)

			if tt.expectedErrors == 0 {
				assert.False(t, result.HasErrors)
				assert.Empty(t, result.Errors)
			} else {
				assert.True(t, result.HasErrors)
				assert.Len(t, result.Errors, tt.expectedErrors)

				for i, expectedType := range tt.expectedTypes {
					assert.Equal(t, expectedType, result.Errors[i].Type)
					assert.NotEmpty(t, result.Errors[i].Details)
				}
			}
		})
	}
}

func TestCalculateCoverage(t *testing.T) {
	tests := []struct {
		name     string
		captions []models.Caption
		start    time.Duration
		end      time.Duration
		expected time.Duration
	}{
		{
			name: "exact match",
			captions: []models.Caption{
				{
					Start: 1 * time.Second,
					End:   5 * time.Second,
					Text:  "Hello",
				},
			},
			start:    1 * time.Second,
			end:      5 * time.Second,
			expected: 4 * time.Second,
		},
		{
			name: "partial overlap - caption starts before range",
			captions: []models.Caption{
				{
					Start: 1 * time.Second,
					End:   7 * time.Second,
					Text:  "Hello",
				},
			},
			start:    3 * time.Second,
			end:      5 * time.Second,
			expected: 2 * time.Second,
		},
		{
			name: "partial overlap - caption ends after range",
			captions: []models.Caption{
				{
					Start: 3 * time.Second,
					End:   7 * time.Second,
					Text:  "Hello",
				},
			},
			start:    1 * time.Second,
			end:      5 * time.Second,
			expected: 2 * time.Second,
		},
		{
			name: "multiple captions with gaps",
			captions: []models.Caption{
				{
					Start: 1 * time.Second,
					End:   3 * time.Second,
					Text:  "Hello",
				},
				{
					Start: 7 * time.Second,
					End:   9 * time.Second,
					Text:  "World",
				},
			},
			start:    0,
			end:      10 * time.Second,
			expected: 4 * time.Second, // 2 + 2 = 4 seconds
		},
		{
			name: "overlapping captions",
			captions: []models.Caption{
				{
					Start: 1 * time.Second,
					End:   5 * time.Second,
					Text:  "Hello",
				},
				{
					Start: 3 * time.Second,
					End:   7 * time.Second,
					Text:  "World",
				},
			},
			start:    0,
			end:      10 * time.Second,
			expected: 8 * time.Second, // 1-7 seconds (4s + 4s = 8s total)
		},
		{
			name:     "no captions",
			captions: []models.Caption{},
			start:    0,
			end:      10 * time.Second,
			expected: 0,
		},
		{
			name: "caption outside range",
			captions: []models.Caption{
				{
					Start: 15 * time.Second,
					End:   20 * time.Second,
					Text:  "Hello",
				},
			},
			start:    0,
			end:      10 * time.Second,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateCoverage(tt.captions, tt.start, tt.end)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatCoverageError(t *testing.T) {
	tests := []struct {
		name           string
		actual         float64
		required       float64
		totalDuration  time.Duration
		covered        time.Duration
		expectedPrefix string
	}{
		{
			name:           "low coverage",
			actual:         50.0,
			required:       80.0,
			totalDuration:  10 * time.Second,
			covered:        5 * time.Second,
			expectedPrefix: "Coverage 50.0% (5.0s of 10.0s) is below required 80.0%.",
		},
		{
			name:           "very low coverage",
			actual:         10.0,
			required:       90.0,
			totalDuration:  100 * time.Second,
			covered:        10 * time.Second,
			expectedPrefix: "Coverage 10.0% (10.0s of 100.0s) is below required 90.0%.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatCoverageError(tt.actual, tt.required, tt.totalDuration, tt.covered)
			assert.Contains(t, result, tt.expectedPrefix)
			assert.Contains(t, result, "Add more captions or extend existing caption durations")
		})
	}
}

func TestMaxDuration(t *testing.T) {
	tests := []struct {
		name     string
		a        time.Duration
		b        time.Duration
		expected time.Duration
	}{
		{"a greater", 5 * time.Second, 3 * time.Second, 5 * time.Second},
		{"b greater", 3 * time.Second, 5 * time.Second, 5 * time.Second},
		{"equal", 3 * time.Second, 3 * time.Second, 3 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maxDuration(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMinDuration(t *testing.T) {
	tests := []struct {
		name     string
		a        time.Duration
		b        time.Duration
		expected time.Duration
	}{
		{"a smaller", 3 * time.Second, 5 * time.Second, 3 * time.Second},
		{"b smaller", 5 * time.Second, 3 * time.Second, 3 * time.Second},
		{"equal", 3 * time.Second, 3 * time.Second, 3 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := minDuration(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}
