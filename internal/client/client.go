package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/Hovhannesmn/ld_proto/pb"
)

// LanguageDetectionService defines the interface for language detection operations
type LanguageDetectionService interface {
	DetectLanguage(ctx context.Context, text string) (*pb.DetectLanguageResponse, error)
	Close() error
}

// LanguageDetectionClient wraps the gRPC client
type LanguageDetectionClient struct {
	client pb.LanguageDetectionServiceClient
	conn   *grpc.ClientConn
}

// NewLanguageDetectionClient creates a new client connection
func NewLanguageDetectionClient(serverAddr string) (*LanguageDetectionClient, error) {
	// Connect to the gRPC server
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}

	client := pb.NewLanguageDetectionServiceClient(conn)

	return &LanguageDetectionClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close closes the client connection
func (c *LanguageDetectionClient) Close() error {
	return c.conn.Close()
}

// DetectLanguage sends a single text string for language detection
func (c *LanguageDetectionClient) DetectLanguage(ctx context.Context, text string) (*pb.DetectLanguageResponse, error) {
	req := &pb.DetectLanguageRequest{
		Text:       text,
		DocumentId: "documentID",
		Metadata: map[string]string{
			"client":    "language-detection-client",
			"timestamp": time.Now().Format(time.RFC3339),
		},
	}

	response, err := c.client.DetectLanguage(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("language detection failed: %w", err)
	}

	return response, nil
}
