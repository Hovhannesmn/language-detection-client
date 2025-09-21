package validator

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"language-detection-client/internal/models"
	pb "github.com/hovman/ld-proto/pb-service/proto"
)

// MockLanguageDetectionService is a mock for the LanguageDetectionService interface
type MockLanguageDetectionService struct{}

func NewMockLanguageDetectionService() *MockLanguageDetectionService {
	return &MockLanguageDetectionService{}
}

func (m *MockLanguageDetectionService) DetectLanguage(ctx context.Context, text string) (*pb.DetectLanguageResponse, error) {
	// Return a mock response - in real implementation, you'd use gomock
	return &pb.DetectLanguageResponse{
		LanguageCode: "en-US",
		Confidence:   0.95,
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

			result := validator.Validate(tt.captions, tt.config)

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

// Integration test that would work with a real service
func TestLanguageValidator_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// This test would require a running language detection service
	// For now, we'll skip it in the test suite
	t.Skip("Integration test requires running language detection service")
}
