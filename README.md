# Language Detection Client

A caption validation client that analyzes subtitle/caption files for language compliance and coverage requirements. This client integrates with a language detection service powered by Amazon Comprehend to ensure captions meet quality standards.

## Features

- **Caption File Validation**: Parse and validate SRT and WebVTT caption files
- **Language Detection**: Uses Amazon Comprehend via gRPC service to detect language and confidence
- **Coverage Analysis**: Validates caption coverage percentage within specified time ranges
- **Duration Filtering**: Only analyzes captions within the specified start/end time range
- **Multi-language Support**: Supports both SRT and WebVTT subtitle formats
- **Docker Ready**: Containerized for easy deployment and testing

## How It Works

### Language Detection Service Integration

This client communicates with a **Language Detection Service** that uses **Amazon Comprehend** to:

1. **Detect Language**: Analyze caption text to identify the language code (e.g., `en-US`, `es-ES`, `fr-FR`)
2. **Calculate Confidence**: Provide a confidence percentage for the language detection
3. **Validate Compliance**: Check if the detected language matches the required language (`en-US`)

### Validation Process

The client performs two main validations:

1. **Language Validation**:
   - Sends all caption text within the specified time range to the language detection service
   - Uses Amazon Comprehend to detect the language and confidence level
   - **Fails validation** if the detected language is NOT `en-US` (English US)
   - Shows error message with detected language and confidence percentage

2. **Coverage Validation**:
   - Calculates what percentage of the specified time range is covered by captions
   - Only considers captions that fall within the `--t_start` and `--t_end` time range
   - **Fails validation** if coverage is below the required percentage (default: 80%)

## Installation

1. **Install dependencies**:
   ```bash
   go mod tidy
   ```

2. **Build the application**:
   ```bash
   go build -o language-detection-client ./cmd/client
   ```

## Usage

### Basic Caption Validation

Validate a caption file with default settings:
```bash
./language-detection-client test-captions/sample.srt
```

### Custom Parameters

Specify time range, coverage requirement, and language detection service:
```bash
./language-detection-client \
  -server=localhost:6011 \
  -t_start=00:00:01 \
  -t_end=00:00:15 \
  -coverage=80 \
  test-captions/sample.srt
```

### Docker Usage

Run with Docker Compose for easy testing:
```bash
make run FILE=test-captions/sample.srt START=00:00:01 END=00:00:15 COVERAGE=80 SERVER=localhost:6011
```

### Command Line Options

- `-server`: Language detection service address (default: `localhost:6011`)
- `-t_start`: Start time for validation range (default: `00:00:00`)
- `-t_end`: End time for validation range (default: `00:40:00`)
- `-coverage`: Required coverage percentage (default: `80`)
- `FILE`: Caption file path (SRT or WebVTT format)

## Example Output

### Successful Validation
When captions pass all validations:
```bash
$ ./language-detection-client test-captions/sample.srt -t_start=00:00:01 -t_end=00:00:20 -coverage=70
# No output - validation passed successfully
```

### Language Validation Error
When detected language is not English (US):
```json
{"type":"incorrect_language","details":"Detected language 'es-ES' (confidence: 98.2%) is not en-US. Expected English (US) captions."}
```

### Coverage Validation Error
When caption coverage is insufficient:
```json
{"type":"caption_coverage","details":"Coverage 65.0% (6.5s of 10.0s) is below required 80.0%. Add more captions or extend existing caption durations to meet coverage requirement."}
```

### Service Connection Error
When language detection service is unavailable:
```json
{"type":"incorrect_language","details":"Language detection failed: language detection failed: rpc error: code = Unavailable desc = connection error: desc = \"transport: Error while dialing: dial tcp 127.0.0.1:6011: connect: connection refused\""}
```

## Prerequisites

- **Language Detection Service**: A gRPC service using Amazon Comprehend must be running and accessible
- **Go 1.19+**: Required for building and running the client
- **Docker** (optional): For containerized deployment and testing

## Language Detection Service Requirements

The language detection service should:

1. **Use Amazon Comprehend**: Integrate with AWS Comprehend for language detection
2. **Return Language Code**: Provide standard language codes (e.g., `en-US`, `es-ES`, `fr-FR`)
3. **Return Confidence**: Include confidence percentage for the detection
4. **Accept gRPC Calls**: Handle `DetectLanguage` gRPC requests
5. **Support Text Analysis**: Process concatenated caption text for analysis

## Testing

### Local Testing

1. **Start the Language Detection Service** (must be running on specified port):
   ```bash
   # Service should be available at localhost:6011 (or your configured port)
   ```

2. **Run validation tests**:
   ```bash
   # Test with English captions (should pass)
   ./language-detection-client test-captions/sample.srt -t_start=00:00:01 -t_end=00:00:20 -coverage=70
   
   # Test with multilingual captions (should show language error)
   ./language-detection-client test-captions/multilingual.srt -t_start=00:00:01 -t_end=00:00:12 -coverage=80
   ```

### Docker Testing

Use the provided Makefile for easy Docker testing:

```bash
# Test with English captions
make run FILE=test-captions/sample.srt START=00:00:01 END=00:00:20 COVERAGE=70 SERVER=localhost:6011

# Test with multilingual captions  
make run FILE=test-captions/multilingual.srt START=00:00:01 END=00:00:12 COVERAGE=80 SERVER=localhost:6011
```

## Supported File Formats

- **SRT (SubRip)**: Standard subtitle format with timestamp and text
- **WebVTT**: Web Video Text Tracks format

## Validation Rules

1. **Language Rule**: All caption text must be detected as English (US) - `en-US`
2. **Coverage Rule**: Captions must cover at least the specified percentage of the time range
3. **Duration Rule**: Only captions within the specified start/end time are analyzed
4. **Format Rule**: File must be valid SRT or WebVTT format
