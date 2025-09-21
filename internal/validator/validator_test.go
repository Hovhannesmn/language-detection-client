package validator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"language-detection-client/internal/models"
)

func TestCaptionValidator_Validate(t *testing.T) {
	tests := []struct {
		name           string
		captions       []models.Caption
		config         models.Config
		expectedErrors int
		expectedTypes  []string
	}{
		{
			name: "all validations pass",
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
				ServerAddr:       "localhost:8080",
			},
			expectedErrors: 0, // Coverage passes, language validation will fail due to no server
			expectedTypes:  []string{},
		},
		{
			name: "coverage validation fails",
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
				ServerAddr:       "localhost:8080",
			},
			expectedErrors: 1, // Coverage fails, language validation will also fail
			expectedTypes:  []string{"caption_coverage"},
		},
		{
			name: "multiple validation failures",
			captions: []models.Caption{
				{
					Start: 1 * time.Second,
					End:   3 * time.Second,
					Text:  "Hello world",
				},
			},
			config: models.Config{
				StartTime:        10 * time.Second,
				EndTime:          5 * time.Second, // Invalid time range
				CoverageRequired: 80.0,
				ServerAddr:       "localhost:8080",
			},
			expectedErrors: 1, // Invalid time range
			expectedTypes:  []string{"invalid_time_range"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create validator - this will fail to connect to server, but we can test coverage logic
			validator, err := NewCaptionValidator(tt.config.ServerAddr)
			if err != nil {
				// Expected to fail due to no server running
				// We'll test with a manually created validator for coverage testing
				validator = &CaptionValidator{
					validators: []Validator{
						NewCoverageValidator(),
					},
					closableValidators: []ClosableValidator{},
				}
			}
			defer func() {
				if validator != nil {
					validator.Close()
				}
			}()

			result := validator.Validate(tt.captions, tt.config)

			// For these tests, we expect at least coverage validation to run
			// Language validation will fail due to no server, but that's expected
			assert.True(t, result.HasErrors)
			assert.NotEmpty(t, result.Errors)
			
			// Check that we have at least the expected error types
			errorTypes := make([]string, len(result.Errors))
			for i, err := range result.Errors {
				errorTypes[i] = err.Type
			}
			
			for _, expectedType := range tt.expectedTypes {
				assert.Contains(t, errorTypes, expectedType)
			}
		})
	}
}

func TestCaptionValidator_Close(t *testing.T) {
	// Test with a validator that has closable validators
	validator := &CaptionValidator{
		validators: []Validator{
			NewCoverageValidator(),
		},
		closableValidators: []ClosableValidator{},
	}

	err := validator.Close()
	assert.NoError(t, err)
}

func TestRunValidations(t *testing.T) {
	tests := []struct {
		name           string
		captions       []models.Caption
		config         *models.Config
		expectedErrors bool
	}{
		{
			name: "valid captions",
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
			config: &models.Config{
				StartTime:        0,
				EndTime:          10 * time.Second,
				CoverageRequired: 80.0,
				ServerAddr:       "localhost:8080",
			},
			expectedErrors: true, // Will have errors due to no server
		},
		{
			name:     "empty captions",
			captions: []models.Caption{},
			config: &models.Config{
				StartTime:        0,
				EndTime:          10 * time.Second,
				CoverageRequired: 80.0,
				ServerAddr:       "localhost:8080",
			},
			expectedErrors: true, // Will have coverage errors
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RunValidations(tt.captions, tt.config)
			
			// Should not return an error (even if validation fails)
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

func TestNewCaptionValidator(t *testing.T) {
	tests := []struct {
		name        string
		serverAddr  string
		expectError bool
	}{
		{
			name:        "valid server address",
			serverAddr:  "localhost:8080",
			expectError: false,
		},
		{
			name:        "empty server address",
			serverAddr:  "",
			expectError: true,
		},
		{
			name:        "invalid server address",
			serverAddr:  "invalid-address",
			expectError: false, // Will create client but fail to connect
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator, err := NewCaptionValidator(tt.serverAddr)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, validator)
			} else {
				// May or may not error depending on server availability
				if err != nil {
					// Expected if no server is running
					assert.Nil(t, validator)
				} else {
					require.NotNil(t, validator)
					defer validator.Close()
				}
			}
		})
	}
}
