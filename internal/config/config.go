package config

import (
	"flag"
	"fmt"
	"language-detection-client/internal/parser"
	"os"
	"strconv"
	"strings"
	"time"

	"language-detection-client/internal/models"
)

// LoadConfig loads configuration from command line flags and arguments
func LoadConfig() (*models.Config, error) {
	// Define command line flags
	tStartStr := flag.String("t_start", "00:00:00", "start time (hh:mm:ss)")
	tEndStr := flag.String("t_end", "00:00:00", "end time (hh:mm:ss)")
	coverageRequired := flag.Float64("coverage", 80.0, "required coverage percentage")
	serverAddr := flag.String("server", "localhost:6011", "gRPC server address")

	// Parse flags
	flag.Parse()

	// Check for required arguments
	if flag.NArg() < 1 {
		return nil, fmt.Errorf("usage: program <flags> captions-filepath")
	}

	filePath := flag.Arg(0)

	// Parse time strings
	tStart, err := parseTime(*tStartStr)
	if err != nil {
		return nil, parser.ErrInvalidTimestamp
	}

	tEnd, err := parseTime(*tEndStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing end time: %v", err)
	}

	// Validate coverage percentage
	if *coverageRequired < 0 || *coverageRequired > 100 {
		return nil, fmt.Errorf("coverage percentage must be between 0 and 100")
	}

	return &models.Config{
		StartTime:        tStart,
		EndTime:          tEnd,
		CoverageRequired: *coverageRequired,
		ServerAddr:       *serverAddr,
		FilePath:         filePath,
	}, nil
}

// parseTime parses a time string in HH:MM:SS format to time.Duration
func parseTime(s string) (time.Duration, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("bad time format, expected HH:MM:SS")
	}

	h, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid hour: %v", err)
	}

	m, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid minute: %v", err)
	}

	sec, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, fmt.Errorf("invalid second: %v", err)
	}

	// Validate ranges
	if h < 0 || h > 23 {
		return 0, fmt.Errorf("hour must be between 0 and 23")
	}
	if m < 0 || m > 59 {
		return 0, fmt.Errorf("minute must be between 0 and 59")
	}
	if sec < 0 || sec > 59 {
		return 0, fmt.Errorf("second must be between 0 and 59")
	}

	return time.Duration(h)*time.Hour + time.Duration(m)*time.Minute + time.Duration(sec)*time.Second, nil
}

// ValidateFileExists checks if the caption file exists and is readable
func ValidateFileExists(filePath string) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("caption file does not exist: %s", filePath)
		}
		return fmt.Errorf("error accessing caption file: %v", err)
	}

	if fileInfo.IsDir() {
		return fmt.Errorf("caption file path is a directory: %s", filePath)
	}

	if fileInfo.Size() == 0 {
		return fmt.Errorf("caption file is empty: %s", filePath)
	}

	return nil
}
