package parser

import "errors"

// Common parser errors
var (
	ErrUnsupportedFormat = errors.New("unsupported caption file format")
	ErrInvalidTimestamp  = errors.New("invalid timestamp format")
)
