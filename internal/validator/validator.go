package validator

import (
	"language-detection-client/internal/models"
)

// Validator interface for different types of validation
type Validator interface {
	Validate(captions []models.Caption, config models.Config) models.ValidationResult
}

// ClosableValidator interface for validators that need cleanup
type ClosableValidator interface {
	Validator
	Close() error
}

// CaptionValidator orchestrates all validation processes
type CaptionValidator struct {
	validators []Validator
	closableValidators []ClosableValidator
}

func NewCaptionValidator(serverAddr string) (*CaptionValidator, error) {
	// Create language validator with server address
	langValidator, err := NewLanguageValidator(serverAddr)
	if err != nil {
		return nil, err
	}

	return &CaptionValidator{
		validators: []Validator{
			NewCoverageValidator(),
			langValidator,
		},
		closableValidators: []ClosableValidator{
			langValidator,
		},
	}, nil
}


func (cv *CaptionValidator) Validate(captions []models.Caption, config models.Config) models.ValidationResult {
	result := models.ValidationResult{
		HasErrors: false,
		Errors:    []models.ValidationError{},
	}

	// Run all validators
	for _, validator := range cv.validators {
		validatorResult := validator.Validate(captions, config)

		// Combine results
		if validatorResult.HasErrors {
			result.HasErrors = true
			result.Errors = append(result.Errors, validatorResult.Errors...)
		}
	}

	return result
}

// Close closes all closable validators
func (cv *CaptionValidator) Close() error {
	for _, validator := range cv.closableValidators {
		if err := validator.Close(); err != nil {
			return err
		}
	}
	return nil
}

// RunValidations runs all configured validators against the captions
func RunValidations(captions []models.Caption, cfg *models.Config) (models.ValidationResult, error) {
	// Create validator with all validation rules
	captionValidator, err := NewCaptionValidator(cfg.ServerAddr)
	if err != nil {
		return models.ValidationResult{}, err
	}
	defer captionValidator.Close()

	return captionValidator.Validate(captions, *cfg), nil
}
