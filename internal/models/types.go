package models

import (
	"time"
)

// Caption represents a single caption entry with timing and text
type Caption struct {
	Start time.Duration // Start time of the caption
	End   time.Duration // End time of the caption
	Text  string        // The caption text content
}


// ValidationError represents a validation error with type and details
type ValidationError struct {
	Type    string `json:"type"`    // Error type (e.g., "caption_coverage", "incorrect_language")
	Details string `json:"details"` // Detailed error description
}


// Config holds all configuration for the validator
type Config struct {
	StartTime        time.Duration // Start time for validation range
	EndTime          time.Duration // End time for validation range
	CoverageRequired float64       // Required coverage percentage
	ServerAddr       string        // Language detection server URL
	FilePath         string        // Path to the caption file
}

// ValidationResult holds the results of validation
type ValidationResult struct {
	HasErrors bool              // Whether validation has errors
	Errors    []ValidationError // List of validation errors
}
