# Language Detection Client

A gRPC client for the Language Detection Service that can send text strings for language analysis.

## Features

- **Single Text Analysis**: Send individual text strings for language detection
- **Batch Processing**: Analyze multiple texts in a single request
- **Health Monitoring**: Check service health status
- **Graceful Shutdown**: Proper signal handling
- **Formatted Output**: Clean, readable results

## Installation

1. **Install dependencies**:
   ```bash
   go mod tidy
   ```

2. **Generate protobuf files**:
   ```bash
   chmod +x scripts/generate_proto.sh
   ./scripts/generate_proto.sh
   ```

## Usage

### Basic Usage

Analyze a single text string:
```bash
go run cmd/client/main.go -text "Hello, world! This is English text."
```

### Custom Server Address

Connect to a different server:
```bash
go run cmd/client/main.go -server "localhost:8080" -text "Hola, mundo!"
```

### Batch Mode

Analyze multiple texts at once:
```bash
go run cmd/client/main.go -batch
```

### Command Line Arguments

Analyze multiple texts from command line:
```bash
go run cmd/client/main.go "Hello world" "Hola mundo" "Bonjour le monde"
```

### Options

- `-server`: gRPC server address (default: `localhost:8090`)
- `-text`: Text to analyze (default: `"Hello, world! This is a test message."`)
- `-doc-id`: Document ID (optional)
- `-batch`: Enable batch mode for testing multiple languages

## Example Output

### Single Detection
```
Language Detection Result:
  Document ID: doc-001
  Detected Language: en-US
  Confidence: 95.50%
  Provider: aws-comprehend
  Processing Time: 150ms
```

### Batch Detection
```
Batch Language Detection Results:
  Total Requests: 5
  Successful: 5
  Failed: 0
  Total Processing Time: 450ms

Result 1:
Language Detection Result:
  Document ID: doc-001
  Detected Language: en-US
  Confidence: 95.50%
  Provider: aws-comprehend
  Processing Time: 120ms

Result 2:
Language Detection Result:
  Document ID: doc-002
  Detected Language: es-ES
  Confidence: 98.20%
  Provider: aws-comprehend
  Processing Time: 110ms
```

## Prerequisites

- The Language Detection gRPC Service must be running
- Go 1.19+ installed
- Protocol Buffers compiler (protoc) installed

## Testing

1. **Start the Language Detection Service** (in another terminal):
   ```bash
   cd ../language-detection-service
   go run cmd/server/main.go
   ```

2. **Run the client**:
   ```bash
   go run cmd/client/main.go
   ```

## Graceful Shutdown

The client supports graceful shutdown via SIGINT (Ctrl+C) or SIGTERM signals. When a shutdown signal is received, ongoing operations are cancelled cleanly.
