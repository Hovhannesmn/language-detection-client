package validator

import (
	"context"
	"fmt"
	"strings"
	"time"

	"language-detection-client/internal/client"
	"language-detection-client/internal/models"
)

// LanguageValidator validates caption language using the gRPC service
type LanguageValidator struct {
	client client.LanguageDetectionService
}

// NewLanguageValidator creates a new language validator
func NewLanguageValidator(serverAddr string) (*LanguageValidator, error) {
	client, err := client.NewLanguageDetectionClient(serverAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create language detection client: %w", err)
	}

	return &LanguageValidator{
		client: client,
	}, nil
}

// Validate checks if the detected language is acceptable (en-US)
func (lv *LanguageValidator) Validate(captions []models.Caption, config models.Config) (models.ValidationResult, error) {
	result := models.ValidationResult{
		HasErrors: false,
		Errors:    []models.ValidationError{},
	}

	// Collect all caption text
	text := collectText(captions)
	if strings.TrimSpace(text) == "" {
		// No text to validate
		return result, nil
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Send text to language detection service
	response, err := lv.client.DetectLanguage(ctx, text)
	if err != nil {
		return result, err
	}

	// Check if language is acceptable (en-US)
	detectedLang := response.LanguageCode
	confidence := response.Confidence

	if detectedLang != "en-US" {
		result.HasErrors = true
		result.Errors = append(result.Errors, models.ValidationError{
			Type:    "incorrect_language",
			Details: fmt.Sprintf("Detected language '%s' (confidence: %.1f%%) is not en-US. Expected English (US) captions.", detectedLang, confidence*100),
		})
	}

	return result, nil
}

// Close closes the language detection client
func (lv *LanguageValidator) Close() error {
	return lv.client.Close()
}

// collectText concatenates all caption text with newlines
func collectText(captions []models.Caption) string {
	var buf strings.Builder
	for _, c := range captions {
		buf.WriteString(c.Text)
		buf.WriteString("\n")
	}
	return buf.String()
}
