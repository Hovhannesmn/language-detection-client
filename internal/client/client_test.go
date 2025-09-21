package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLanguageDetectionClient_DetectLanguage(t *testing.T) {
	// Basic test without mocking - just test that the method exists and handles nil client
	client := &LanguageDetectionClient{
		client: nil,
		conn:   nil,
	}

	ctx := context.Background()
	result, err := client.DetectLanguage(ctx, "test text")

	// Should error due to nil client
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestLanguageDetectionClient_Close(t *testing.T) {
	client := &LanguageDetectionClient{
		client: nil,
		conn:   nil, // In real scenario, this would be a gRPC connection
	}

	// For testing, we can't easily mock the connection close
	// In a real test, you'd need to create a mock connection or use dependency injection
	err := client.Close()
	
	// This should not panic even with nil connection
	assert.NoError(t, err)
}

func TestNewLanguageDetectionClient(t *testing.T) {
	tests := []struct {
		name        string
		serverAddr  string
		expectError bool
	}{
		{
			name:        "valid server address",
			serverAddr:  "localhost:8080",
			expectError: false,
		},
		{
			name:        "empty server address",
			serverAddr:  "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: This test will actually try to connect to the server
			// In a real test environment, you might want to mock the gRPC dial
			if tt.expectError {
				client, err := NewLanguageDetectionClient(tt.serverAddr)
				assert.Error(t, err)
				assert.Nil(t, client)
			}
			// For valid addresses, we can't easily test without a real server
		})
	}
}
