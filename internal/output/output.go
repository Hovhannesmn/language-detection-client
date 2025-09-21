package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"language-detection-client/internal/models"
)

// OutputHandler handles output formatting and writing
type OutputHandler struct {
	resultWriter io.Writer // For validation results (JSON)
	logWriter    io.Writer // For logs and errors
}

// NewOutputHandler creates a new output handler
func NewOutputHandler() *OutputHandler {
	return &OutputHandler{
		resultWriter: os.Stdout, // Results go to stdout
		logWriter:    os.Stderr, // Logs go to stderr
	}
}


// WriteValidationErrors writes validation errors to the result stream
func (oh *OutputHandler) WriteValidationErrors(errors []models.ValidationError) error {
	for _, err := range errors {
		if err := oh.WriteValidationError(err); err != nil {
			return fmt.Errorf("failed to write validation error: %v", err)
		}
	}
	return nil
}

// WriteValidationError writes a single validation error to the result stream
func (oh *OutputHandler) WriteValidationError(err models.ValidationError) error {
	jsonData, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		return fmt.Errorf("error marshaling validation error: %v", marshalErr)
	}

	_, writeErr := fmt.Fprintln(oh.resultWriter, string(jsonData))
	if writeErr != nil {
		return fmt.Errorf("error writing validation result: %v", writeErr)
	}

	return nil
}

// WriteError writes an error message to the log stream
func (oh *OutputHandler) WriteError(format string, args ...interface{}) {
	fmt.Fprintf(oh.logWriter, "ERROR: "+format, args...)
}

// OutputValidationResults outputs the validation results using the output handler
func (oh *OutputHandler) OutputValidationResults(result models.ValidationResult) error {
	if result.HasErrors {
		// Write each validation error as JSON to stdout
		return oh.WriteValidationErrors(result.Errors)
	}
	// If no errors, no output (as per requirements)
	return nil
}
