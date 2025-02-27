package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/MateSousa/overtime-script/internal/config"
	"github.com/MateSousa/overtime-script/pkg/adapters/exporters"
	"github.com/MateSousa/overtime-script/pkg/adapters/notification"
	"github.com/MateSousa/overtime-script/pkg/adapters/repositories"
	"github.com/MateSousa/overtime-script/pkg/domain/usecases"
	"github.com/MateSousa/overtime-script/pkg/infrastructure/kubernetes"
)

func main() {
	// Create context
	ctx := context.Background()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Create Kubernetes client
	k8sClient, err := kubernetes.NewInClusterClient()
	if err != nil {
		fmt.Printf("Error creating Kubernetes client: %v\n", err)
		os.Exit(1)
	}

	// Create repositories and services
	overtimeRepo := repositories.NewKubernetesOvertimeRepository(k8sClient, cfg.Namespace)
	excelExporter := exporters.NewExcelReportExporter()
	emailService := notification.NewSESEmailService(
		cfg.SenderEmail,
		cfg.RecipientEmail,
		cfg.AWSRegion,
	)

	// Create use case
	overtimeUseCase := usecases.NewOvertimeUseCase(
		overtimeRepo,
		excelExporter,
		emailService,
	)

	// Handle testing mode or normal operation
	if cfg.TestingMode {
		fmt.Println("Running in test mode...")
		if err := overtimeUseCase.TestMonthlyReport(ctx); err != nil {
			fmt.Printf("Error in test mode: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Test completed successfully!")
		os.Exit(0)
	}

	// Check if it's the first day of the month
	now := time.Now().UTC()
	if now.Day() == 1 {
		fmt.Println("First day of the month, generating monthly report...")
		if err := overtimeUseCase.GenerateMonthlyReport(ctx); err != nil {
			fmt.Printf("Error generating monthly report: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Monthly report sent successfully!")
	} else {
		fmt.Println("Processing yesterday's overtime entries...")
		if err := overtimeUseCase.ProcessYesterdayOvertime(ctx); err != nil {
			fmt.Printf("Error processing yesterday's overtime: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Yesterday's overtime processed successfully!")
	}
}