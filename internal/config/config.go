package config

import (
	"fmt"
	"os"
)

// Config holds application configuration
type Config struct {
	// Kubernetes configuration
	Namespace string

	// Email configuration
	SenderEmail   string
	RecipientEmail string
	AWSRegion     string

	// Application mode
	TestingMode bool
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load namespace, defaulting to "default"
	namespace := os.Getenv("NAMESPACE")
	if namespace == "" {
		namespace = "default"
	}

	// Load email configuration
	senderEmail := os.Getenv("SENDER_EMAIL")
	if senderEmail == "" {
		return nil, fmt.Errorf("SENDER_EMAIL environment variable is required")
	}

	recipientEmail := os.Getenv("RECIPIENT_EMAIL")
	if recipientEmail == "" {
		return nil, fmt.Errorf("RECIPIENT_EMAIL environment variable is required")
	}

	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		return nil, fmt.Errorf("AWS_REGION environment variable is required")
	}

	// Check if we're in testing mode
	testingMode := os.Getenv("TESTING") == "true"

	return &Config{
		Namespace:     namespace,
		SenderEmail:   senderEmail,
		RecipientEmail: recipientEmail,
		AWSRegion:     awsRegion,
		TestingMode:   testingMode,
	}, nil
}