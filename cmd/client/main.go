package main

import (
	"language-detection-client/internal/config"
	"language-detection-client/internal/output"
	"language-detection-client/internal/parser"
	"language-detection-client/internal/validator"
	"os"
)

func main() {
	outputHandler := output.NewOutputHandler()

	// CLI flags
	// Load configuration from command line arguments
	cfg, err := config.LoadConfig()
	if err != nil {
		outputHandler.WriteError("%v\n", err)
		os.Exit(1)
	}

	// Validate that the caption file exists and is accessible
	if err := config.ValidateFileExists(cfg.FilePath); err != nil {
		outputHandler.WriteError("%v\n", err)
		os.Exit(1)
	}

	// Read caption file content
	content, err := parser.ReadCaptionFile(cfg.FilePath)
	if err != nil {
		outputHandler.WriteError("%v\n", err)
		os.Exit(1)
	}

	// Parse caption file using appropriate parser
	captions, err := parser.ParseCaptionFile(content)
	if err != nil {
		outputHandler.WriteError("%v\n", err)
		os.Exit(1)
	}

	// Run all validations
	validationResult, err := validator.RunValidations(captions, cfg)
	if err != nil {
		outputHandler.WriteError("Failed to run validations: %v\n", err)
		os.Exit(1)
	}

	// Output validation results
	if err := outputHandler.OutputValidationResults(validationResult); err != nil {
		outputHandler.WriteError("Failed to output results: %v\n", err)
		os.Exit(1)
	}
}

