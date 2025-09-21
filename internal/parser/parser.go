package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"language-detection-client/internal/models"
)

// SRTRegex matches SRT timestamp format: HH:MM:SS,mmm --> HH:MM:SS,mmm
var SRTRegex = regexp.MustCompile(`(\d{2}):(\d{2}):(\d{2}),(\d{3}) --> (\d{2}):(\d{2}):(\d{2}),(\d{3})`)

// VTTRegex matches WebVTT timestamp format: HH:MM:SS.mmm --> HH:MM:SS.mmm
var VTTRegex = regexp.MustCompile(`(\d{2}):(\d{2}):(\d{2})\.(\d{3}) --> (\d{2}):(\d{2}):(\d{2})\.(\d{3})`)

// CaptionParser interface for different caption formats
type CaptionParser interface {
	Parse(content string) ([]models.Caption, error)
	IsSupported(content string) bool
}

// SRTParser handles SRT format parsing
type SRTParser struct{}

// NewSRTParser creates a new SRT parser
func NewSRTParser() *SRTParser {
	return &SRTParser{}
}

// IsSupported checks if the content is in SRT format
func (p *SRTParser) IsSupported(content string) bool {
	return SRTRegex.MatchString(content)
}

// Parse parses SRT content and returns captions
func (p *SRTParser) Parse(content string) ([]models.Caption, error) {
	var captions []models.Caption
	scanner := bufio.NewScanner(strings.NewReader(content))
	var textBuffer []string
	var start, end time.Duration

	for scanner.Scan() {
		line := scanner.Text()
		if match := SRTRegex.FindStringSubmatch(line); len(match) > 0 {
			start = parseDuration(match[1], match[2], match[3], match[4])
			end = parseDuration(match[5], match[6], match[7], match[8])
			textBuffer = []string{}
		} else if line == "" {
			if len(textBuffer) > 0 {
				captions = append(captions, models.Caption{
					Start: start,
					End:   end,
					Text:  strings.Join(textBuffer, " "),
				})
			}
		} else if !isNumeric(line) {
			textBuffer = append(textBuffer, line)
		}
	}

	// Handle last caption if file doesn't end with empty line
	if len(textBuffer) > 0 {
		captions = append(captions, models.Caption{
			Start: start,
			End:   end,
			Text:  strings.Join(textBuffer, " "),
		})
	}

	return captions, nil
}

// VTTParser handles WebVTT format parsing
type VTTParser struct{}

// NewVTTParser creates a new WebVTT parser
func NewVTTParser() *VTTParser {
	return &VTTParser{}
}

// IsSupported checks if the content is in WebVTT format
func (p *VTTParser) IsSupported(content string) bool {
	return strings.HasPrefix(content, "WEBVTT")
}

// Parse parses WebVTT content and returns captions
func (p *VTTParser) Parse(content string) ([]models.Caption, error) {
	var captions []models.Caption
	lines := strings.Split(content, "\n")
	var textBuffer []string
	var start, end time.Duration
	headerSkipped := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and header
		if line == "" {
			if len(textBuffer) > 0 {
				captions = append(captions, models.Caption{
					Start: start,
					End:   end,
					Text:  strings.Join(textBuffer, " "),
				})
				textBuffer = []string{}
			}
			continue
		}

		if !headerSkipped {
			if strings.HasPrefix(line, "WEBVTT") || strings.HasPrefix(line, "NOTE") {
				headerSkipped = true
				continue
			}
			headerSkipped = true
		}

		if match := VTTRegex.FindStringSubmatch(line); len(match) > 0 {
			start = parseDuration(match[1], match[2], match[3], match[4])
			end = parseDuration(match[5], match[6], match[7], match[8])
			textBuffer = []string{}
		} else if !isNumeric(line) && !strings.HasPrefix(line, "NOTE") {
			// Only add non-numeric, non-note lines as text
			textBuffer = append(textBuffer, line)
		}
	}

	// Handle last caption if file doesn't end with empty line
	if len(textBuffer) > 0 {
		captions = append(captions, models.Caption{
			Start: start,
			End:   end,
			Text:  strings.Join(textBuffer, " "),
		})
	}

	return captions, nil
}

// parseDuration converts hour, minute, second, and millisecond strings to time.Duration
func parseDuration(h, m, s, ms string) time.Duration {
	hi, err1 := strconv.Atoi(h)
	mi, err2 := strconv.Atoi(m)
	si, err3 := strconv.Atoi(s)
	msi, err4 := strconv.Atoi(ms)

	// If any conversion fails, return 0 (should not happen with valid regex)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return 0
	}

	return time.Duration(hi)*time.Hour + time.Duration(mi)*time.Minute +
		time.Duration(si)*time.Second + time.Duration(msi)*time.Millisecond
}

// isNumeric checks if a string represents a number
func isNumeric(s string) bool {
	_, err := strconv.Atoi(strings.TrimSpace(s))
	return err == nil
}

// GetParser returns the appropriate parser for the given content
func GetParser(content string) (CaptionParser, error) {
	// Try WebVTT first (more specific check)
	if strings.HasPrefix(content, "WEBVTT") {
		return NewVTTParser(), nil
	}

	// Try SRT
	if SRTRegex.MatchString(content) {
		return NewSRTParser(), nil
	}

	return nil, ErrUnsupportedFormat
}

// ReadCaptionFile reads the content of the caption file
func ReadCaptionFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading file: %v", err)
	}
	return string(data), nil
}

// ParseCaptionFile parses the caption file content using the appropriate parser
func ParseCaptionFile(content string) ([]models.Caption, error) {
	// Get the appropriate parser for the content
	captionParser, err := GetParser(content)
	if err != nil {
		return nil, fmt.Errorf("parse error: %v", err)
	}

	// Parse the content
	captions, err := captionParser.Parse(content)
	if err != nil {
		return nil, fmt.Errorf("parse error: %v", err)
	}

	return captions, nil
}
