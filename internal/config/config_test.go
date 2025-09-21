package config

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"language-detection-client/internal/models"
)

func TestLoadConfig(t *testing.T) {
	// Save original args and restore after test
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	}()

	tests := []struct {
		name        string
		args        []string
		expected    *models.Config
		expectError bool
	}{
		{
			name: "valid configuration with all flags",
			args: []string{
				"program",
				"-t_start=00:01:30",
				"-t_end=00:02:45",
				"-coverage=85.5",
				"-server=localhost:8080",
				"test.srt",
			},
			expected: &models.Config{
				StartTime:        1*time.Minute + 30*time.Second,
				EndTime:          2*time.Minute + 45*time.Second,
				CoverageRequired: 85.5,
				ServerAddr:       "localhost:8080",
				FilePath:         "test.srt",
			},
			expectError: false,
		},
		{
			name: "minimal configuration with defaults",
			args: []string{
				"program",
				"test.srt",
			},
			expected: &models.Config{
				StartTime:        0,
				EndTime:          0,
				CoverageRequired: 80.0,
				ServerAddr:       "localhost:6011",
				FilePath:         "test.srt",
			},
			expectError: false,
		},
		{
			name: "missing file path",
			args: []string{
				"program",
				"-t_start=00:01:00",
			},
			expected:    nil,
			expectError: true,
		},
		{
			name: "invalid start time format",
			args: []string{
				"program",
				"-t_start=invalid",
				"test.srt",
			},
			expected:    nil,
			expectError: true,
		},
		{
			name: "invalid end time format",
			args: []string{
				"program",
				"-t_end=25:70:80",
				"test.srt",
			},
			expected:    nil,
			expectError: true,
		},
		{
			name: "coverage too high",
			args: []string{
				"program",
				"-coverage=150",
				"test.srt",
			},
			expected:    nil,
			expectError: true,
		},
		{
			name: "coverage too low",
			args: []string{
				"program",
				"-coverage=-10",
				"test.srt",
			},
			expected:    nil,
			expectError: true,
		},
		{
			name: "valid edge case coverage",
			args: []string{
				"program",
				"-coverage=0",
				"test.srt",
			},
			expected: &models.Config{
				StartTime:        0,
				EndTime:          0,
				CoverageRequired: 0,
				ServerAddr:       "localhost:6011",
				FilePath:         "test.srt",
			},
			expectError: false,
		},
		{
			name: "maximum valid coverage",
			args: []string{
				"program",
				"-coverage=100",
				"test.srt",
			},
			expected: &models.Config{
				StartTime:        0,
				EndTime:          0,
				CoverageRequired: 100,
				ServerAddr:       "localhost:6011",
				FilePath:         "test.srt",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flag.CommandLine for each test
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			os.Args = tt.args

			result, err := LoadConfig()

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.expected.StartTime, result.StartTime)
				assert.Equal(t, tt.expected.EndTime, result.EndTime)
				assert.Equal(t, tt.expected.CoverageRequired, result.CoverageRequired)
				assert.Equal(t, tt.expected.ServerAddr, result.ServerAddr)
				assert.Equal(t, tt.expected.FilePath, result.FilePath)
			}
		})
	}
}

func TestParseTime(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    time.Duration
		expectError bool
	}{
		{
			name:        "valid time",
			input:       "01:30:45",
			expected:    1*time.Hour + 30*time.Minute + 45*time.Second,
			expectError: false,
		},
		{
			name:        "midnight",
			input:       "00:00:00",
			expected:    0,
			expectError: false,
		},
		{
			name:        "maximum valid time",
			input:       "23:59:59",
			expected:    23*time.Hour + 59*time.Minute + 59*time.Second,
			expectError: false,
		},
		{
			name:        "invalid format - too few parts",
			input:       "01:30",
			expected:    0,
			expectError: true,
		},
		{
			name:        "invalid format - too many parts",
			input:       "01:30:45:123",
			expected:    0,
			expectError: true,
		},
		{
			name:        "invalid hour - negative",
			input:       "-1:30:45",
			expected:    0,
			expectError: true,
		},
		{
			name:        "invalid hour - too high",
			input:       "24:30:45",
			expected:    0,
			expectError: true,
		},
		{
			name:        "invalid minute - negative",
			input:       "01:-30:45",
			expected:    0,
			expectError: true,
		},
		{
			name:        "invalid minute - too high",
			input:       "01:60:45",
			expected:    0,
			expectError: true,
		},
		{
			name:        "invalid second - negative",
			input:       "01:30:-45",
			expected:    0,
			expectError: true,
		},
		{
			name:        "invalid second - too high",
			input:       "01:30:60",
			expected:    0,
			expectError: true,
		},
		{
			name:        "non-numeric hour",
			input:       "aa:30:45",
			expected:    0,
			expectError: true,
		},
		{
			name:        "non-numeric minute",
			input:       "01:bb:45",
			expected:    0,
			expectError: true,
		},
		{
			name:        "non-numeric second",
			input:       "01:30:cc",
			expected:    0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseTime(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, time.Duration(0), result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestValidateFileExists(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "test_caption_*.srt")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// Write some content to the file
	_, err = tempFile.WriteString("1\n00:00:01,000 --> 00:00:04,000\nHello world")
	require.NoError(t, err)
	tempFile.Close()

	tests := []struct {
		name        string
		filePath    string
		expectError bool
	}{
		{
			name:        "existing file",
			filePath:    tempFile.Name(),
			expectError: false,
		},
		{
			name:        "non-existent file",
			filePath:    "nonexistent.srt",
			expectError: true,
		},
		{
			name:        "empty file path",
			filePath:    "",
			expectError: true,
		},
		{
			name:        "directory path",
			filePath:    "/tmp",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFileExists(tt.filePath)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	// Test empty file
	tempEmptyFile, err := os.CreateTemp("", "empty_*.srt")
	require.NoError(t, err)
	defer os.Remove(tempEmptyFile.Name())
	tempEmptyFile.Close()

	err = ValidateFileExists(tempEmptyFile.Name())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty")
}
