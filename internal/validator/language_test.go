package validator

import (
	"context"
	"fmt"
	"testing"
	"time"

	pb "github.com/Hovhannesmn/ld_proto/pb"
	"github.com/stretchr/testify/assert"
	"language-detection-client/internal/models"
)

// MockLanguageDetectionService is a mock for the LanguageDetectionService interface
type MockLanguageDetectionService struct {
	shouldReturnError bool
	returnLanguage    string
	returnConfidence  float32
}

func NewMockLanguageDetectionService() *MockLanguageDetectionService {
	return &MockLanguageDetectionService{
		returnLanguage:   "en-US",
		returnConfidence: 0.95,
	}
}

func (m *MockLanguageDetectionService) DetectLanguage(ctx context.Context, text string) (*pb.DetectLanguageResponse, error) {
	if m.shouldReturnError {
		return nil, fmt.Errorf("mock connection error")
	}

	return &pb.DetectLanguageResponse{
		LanguageCode: m.returnLanguage,
		Confidence:   m.returnConfidence,
	}, nil
}

func (m *MockLanguageDetectionService) Close() error {
	return nil
}

func TestLanguageValidator_Validate(t *testing.T) {
	tests := []struct {
		name           string
		captions       []models.Caption
		config         models.Config
		expectedErrors bool
	}{
		{
			name:     "empty captions",
			captions: []models.Caption{},
			config: models.Config{
				StartTime:  0,
				EndTime:    10 * time.Second,
				ServerAddr: "localhost:8080",
			},
			expectedErrors: false,
		},
		{
			name: "captions with empty text",
			captions: []models.Caption{
				{
					Start: 1 * time.Second,
					End:   3 * time.Second,
					Text:  "",
				},
				{
					Start: 4 * time.Second,
					End:   6 * time.Second,
					Text:  "   ", // whitespace only
				},
			},
			config: models.Config{
				StartTime:  0,
				EndTime:    10 * time.Second,
				ServerAddr: "localhost:8080",
			},
			expectedErrors: false,
		},
		{
			name: "captions with text - will pass with mock",
			captions: []models.Caption{
				{
					Start: 1 * time.Second,
					End:   3 * time.Second,
					Text:  "Hello world",
				},
			},
			config: models.Config{
				StartTime:  0,
				EndTime:    10 * time.Second,
				ServerAddr: "localhost:8080",
			},
			expectedErrors: false, // Will pass with mock returning en-US
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock service
			mockService := NewMockLanguageDetectionService()

			// Create validator with mock service
			validator := &LanguageValidator{
				client: mockService,
			}

			result, err := validator.Validate(tt.captions, tt.config)
			assert.NoError(t, err)

			if tt.expectedErrors {
				assert.True(t, result.HasErrors)
				assert.NotEmpty(t, result.Errors)
			} else {
				assert.False(t, result.HasErrors)
				assert.Empty(t, result.Errors)
			}
		})
	}
}

func TestLanguageValidator_Validate_WithError(t *testing.T) {
	// Create a mock service that returns an error
	mockService := &MockLanguageDetectionService{
		shouldReturnError: true,
	}

	// Create validator with mock service
	validator := &LanguageValidator{
		client: mockService,
	}

	captions := []models.Caption{
		{
			Start: 1 * time.Second,
			End:   3 * time.Second,
			Text:  "Hello world",
		},
	}

	config := models.Config{
		StartTime:  0,
		EndTime:    10 * time.Second,
		ServerAddr: "localhost:8080",
	}

	// This should return an error now
	result, err := validator.Validate(captions, config)

	// Should return an error
	assert.Error(t, err)
	assert.Equal(t, "mock connection error", err.Error())

	// Result should be empty/default
	assert.False(t, result.HasErrors)
	assert.Empty(t, result.Errors)
}

func TestCollectText(t *testing.T) {
	tests := []struct {
		name     string
		captions []models.Caption
		expected string
	}{
		{
			name:     "empty captions",
			captions: []models.Caption{},
			expected: "",
		},
		{
			name: "single caption",
			captions: []models.Caption{
				{
					Start: 1 * time.Second,
					End:   3 * time.Second,
					Text:  "Hello world",
				},
			},
			expected: "Hello world\n",
		},
		{
			name: "multiple captions",
			captions: []models.Caption{
				{
					Start: 1 * time.Second,
					End:   3 * time.Second,
					Text:  "Hello world",
				},
				{
					Start: 4 * time.Second,
					End:   6 * time.Second,
					Text:  "Goodbye world",
				},
			},
			expected: "Hello world\nGoodbye world\n",
		},
		{
			name: "captions with empty text",
			captions: []models.Caption{
				{
					Start: 1 * time.Second,
					End:   3 * time.Second,
					Text:  "Hello world",
				},
				{
					Start: 4 * time.Second,
					End:   6 * time.Second,
					Text:  "",
				},
				{
					Start: 7 * time.Second,
					End:   9 * time.Second,
					Text:  "Goodbye world",
				},
			},
			expected: "Hello world\n\nGoodbye world\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := collectText(tt.captions)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLanguageValidator_Close(t *testing.T) {
	validator := &LanguageValidator{
		client: NewMockLanguageDetectionService(),
	}

	err := validator.Close()
	assert.NoError(t, err)
}
