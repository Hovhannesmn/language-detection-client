package validator

import (
	"fmt"
	"time"

	"language-detection-client/internal/models"
)

// CoverageValidator validates caption coverage percentage
type CoverageValidator struct{}

func NewCoverageValidator() *CoverageValidator {
	return &CoverageValidator{}
}

func (cv *CoverageValidator) Validate(captions []models.Caption, config models.Config) (models.ValidationResult, error) {
	result := models.ValidationResult{
		HasErrors: false,
		Errors:    []models.ValidationError{},
	}

	// Calculate total duration
	totalDuration := config.EndTime - config.StartTime
	if totalDuration <= 0 {
		result.HasErrors = true
		result.Errors = append(result.Errors, models.ValidationError{
			Type:    "invalid_time_range",
			Details: "End time must be after start time",
		})
		return result, nil
	}

	// Calculate covered duration
	covered := calculateCoverage(captions, config.StartTime, config.EndTime)
	percent := (covered.Seconds() / totalDuration.Seconds()) * 100

	// Check if coverage meets requirement
	if percent < config.CoverageRequired {
		result.HasErrors = true
		result.Errors = append(result.Errors, models.ValidationError{
			Type:    "caption_coverage",
			Details: formatCoverageError(percent, config.CoverageRequired, totalDuration, covered),
		})
	}

	return result, nil
}

// calculateCoverage calculates the total covered duration within the specified time range
func calculateCoverage(captions []models.Caption, start, end time.Duration) time.Duration {
	var covered time.Duration
	for _, c := range captions {
		// Skip captions that don't overlap with the time range
		if c.End < start || c.Start > end {
			continue
		}

		// Calculate the overlapping portion
		s := maxDuration(c.Start, start)
		e := minDuration(c.End, end)
		covered += e - s
	}
	return covered
}

// formatCoverageError formats the coverage error message
func formatCoverageError(actual, required float64, totalDuration, covered time.Duration) string {
	return fmt.Sprintf("Coverage %.1f%% (%.1fs of %.1fs) is below required %.1f%%. Add more captions or extend existing caption durations to meet coverage requirement.",
		actual, covered.Seconds(), totalDuration.Seconds(), required)
}

// maxDuration returns the maximum of two durations
func maxDuration(a, b time.Duration) time.Duration {
	if a > b {
		return a
	}
	return b
}

// minDuration returns the minimum of two durations
func minDuration(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
