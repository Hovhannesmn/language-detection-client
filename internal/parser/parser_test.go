package parser

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"language-detection-client/internal/models"
)

func TestSRTParser_IsSupported(t *testing.T) {
	parser := NewSRTParser()

	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name: "valid SRT format",
			content: `1
00:00:01,000 --> 00:00:04,000
Hello world

2
00:00:05,000 --> 00:00:08,000
Goodbye world`,
			expected: true,
		},
		{
			name:     "empty content",
			content:  "",
			expected: false,
		},
		{
			name:     "invalid format",
			content:  "Just plain text",
			expected: false,
		},
		{
			name: "WebVTT format",
			content: `WEBVTT

1
00:00:01.000 --> 00:00:04.000
Hello world`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.IsSupported(tt.content)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSRTParser_Parse(t *testing.T) {
	parser := NewSRTParser()

	tests := []struct {
		name          string
		content       string
		expected      []models.Caption
		expectedError bool
	}{
		{
			name: "valid SRT with single caption",
			content: `1
00:00:01,000 --> 00:00:04,000
Hello world`,
			expected: []models.Caption{
				{
					Start: 1*time.Second + 0*time.Millisecond,
					End:   4*time.Second + 0*time.Millisecond,
					Text:  "Hello world",
				},
			},
			expectedError: false,
		},
		{
			name: "valid SRT with multiple captions",
			content: `1
00:00:01,000 --> 00:00:04,000
Hello world

2
00:00:05,000 --> 00:00:08,000
Goodbye world

3
00:00:10,500 --> 00:00:13,750
Multi-line
caption text`,
			expected: []models.Caption{
				{
					Start: 1*time.Second + 0*time.Millisecond,
					End:   4*time.Second + 0*time.Millisecond,
					Text:  "Hello world",
				},
				{
					Start: 5*time.Second + 0*time.Millisecond,
					End:   8*time.Second + 0*time.Millisecond,
					Text:  "Goodbye world",
				},
				{
					Start: 10*time.Second + 500*time.Millisecond,
					End:   13*time.Second + 750*time.Millisecond,
					Text:  "Multi-line caption text",
				},
			},
			expectedError: false,
		},
		{
			name: "SRT with empty lines and numbers",
			content: `1
00:00:01,000 --> 00:00:04,000
Hello world

2
00:00:05,000 --> 00:00:08,000
123
Goodbye world`,
			expected: []models.Caption{
				{
					Start: 1*time.Second + 0*time.Millisecond,
					End:   4*time.Second + 0*time.Millisecond,
					Text:  "Hello world",
				},
				{
					Start: 5*time.Second + 0*time.Millisecond,
					End:   8*time.Second + 0*time.Millisecond,
					Text:  "Goodbye world", // The parser skips numeric lines
				},
			},
			expectedError: false,
		},
		{
			name: "empty content",
			content: "",
			expected: []models.Caption{},
			expectedError: false,
		},
		{
			name: "SRT without timestamps",
			content: `1
Hello world
2
Goodbye world`,
			expected: []models.Caption{
				{
					Start: 0, // No timestamp found, so start/end are 0
					End:   0,
					Text:  "Hello world Goodbye world",
				},
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.Parse(tt.content)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				require.Len(t, result, len(tt.expected))
				
				for i, caption := range result {
					assert.Equal(t, tt.expected[i].Start, caption.Start, "Start time mismatch for caption %d", i)
					assert.Equal(t, tt.expected[i].End, caption.End, "End time mismatch for caption %d", i)
					assert.Equal(t, tt.expected[i].Text, caption.Text, "Text mismatch for caption %d", i)
				}
			}
		})
	}
}

func TestVTTParser_IsSupported(t *testing.T) {
	parser := NewVTTParser()

	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name: "valid WebVTT format",
			content: `WEBVTT

1
00:00:01.000 --> 00:00:04.000
Hello world`,
			expected: true,
		},
		{
			name: "WebVTT with notes",
			content: `WEBVTT
NOTE This is a note

1
00:00:01.000 --> 00:00:04.000
Hello world`,
			expected: true,
		},
		{
			name:     "SRT format",
			content:  "1\n00:00:01,000 --> 00:00:04,000\nHello world",
			expected: false,
		},
		{
			name:     "empty content",
			content:  "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.IsSupported(tt.content)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestVTTParser_Parse(t *testing.T) {
	parser := NewVTTParser()

	tests := []struct {
		name          string
		content       string
		expected      []models.Caption
		expectedError bool
	}{
		{
			name: "valid WebVTT with single caption",
			content: `WEBVTT

1
00:00:01.000 --> 00:00:04.000
Hello world`,
			expected: []models.Caption{
				{
					Start: 1*time.Second + 0*time.Millisecond,
					End:   4*time.Second + 0*time.Millisecond,
					Text:  "Hello world",
				},
			},
			expectedError: false,
		},
		{
			name: "WebVTT with notes and multiple captions",
			content: `WEBVTT
NOTE This is a note

1
00:00:01.000 --> 00:00:04.000
Hello world

2
00:00:05.000 --> 00:00:08.000
Goodbye world
NOTE Another note`,
			expected: []models.Caption{
				{
					Start: 1*time.Second + 0*time.Millisecond,
					End:   4*time.Second + 0*time.Millisecond,
					Text:  "Hello world",
				},
				{
					Start: 5*time.Second + 0*time.Millisecond,
					End:   8*time.Second + 0*time.Millisecond,
					Text:  "Goodbye world",
				},
			},
			expectedError: false,
		},
		{
			name: "empty WebVTT",
			content: `WEBVTT`,
			expected: []models.Caption{},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.Parse(tt.content)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				require.Len(t, result, len(tt.expected))
				
				for i, caption := range result {
					assert.Equal(t, tt.expected[i].Start, caption.Start, "Start time mismatch for caption %d", i)
					assert.Equal(t, tt.expected[i].End, caption.End, "End time mismatch for caption %d", i)
					assert.Equal(t, tt.expected[i].Text, caption.Text, "Text mismatch for caption %d", i)
				}
			}
		})
	}
}

func TestGetParser(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		expectedType  string
		expectedError bool
	}{
		{
			name:          "WebVTT content",
			content:       "WEBVTT\n\n1\n00:00:01.000 --> 00:00:04.000\nHello world",
			expectedType:  "*parser.VTTParser",
			expectedError: false,
		},
		{
			name:          "SRT content",
			content:       "1\n00:00:01,000 --> 00:00:04,000\nHello world",
			expectedType:  "*parser.SRTParser",
			expectedError: false,
		},
		{
			name:          "unsupported format",
			content:       "Just plain text without timestamps",
			expectedType:  "",
			expectedError: true,
		},
		{
			name:          "empty content",
			content:       "",
			expectedType:  "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, err := GetParser(tt.content)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, parser)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, parser)
				assert.Equal(t, tt.expectedType, fmt.Sprintf("%T", parser))
			}
		})
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name     string
		h        string
		m        string
		s        string
		ms       string
		expected time.Duration
	}{
		{
			name:     "basic duration",
			h:        "0", m: "1", s: "30", ms: "500",
			expected: 1*time.Minute + 30*time.Second + 500*time.Millisecond,
		},
		{
			name:     "hour duration",
			h:        "1", m: "0", s: "0", ms: "0",
			expected: 1 * time.Hour,
		},
		{
			name:     "zero duration",
			h:        "0", m: "0", s: "0", ms: "0",
			expected: 0,
		},
		{
			name:     "invalid hour",
			h:        "invalid", m: "1", s: "30", ms: "500",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseDuration(tt.h, tt.m, tt.s, tt.ms)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsNumeric(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"positive number", "123", true},
		{"zero", "0", true},
		{"negative number", "-123", true}, // Our function treats negative as valid number
		{"decimal", "123.45", false},
		{"text", "hello", false},
		{"empty", "", false},
		{"spaces", "  123  ", true},
		{"mixed", "123abc", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isNumeric(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
