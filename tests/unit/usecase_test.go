package unit

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MateSousa/overtime-script/pkg/domain/entities"
	"github.com/MateSousa/overtime-script/pkg/domain/usecases"
	"github.com/MateSousa/overtime-script/tests/unit/mocks"
)

func TestProcessYesterdayOvertime(t *testing.T) {
	// Create mocks
	repo := mocks.NewMockOvertimeRepository()
	exporter := mocks.NewMockReportExporter()
	notifier := mocks.NewMockNotificationService()
	
	// Create use case
	uc := usecases.NewOvertimeUseCase(repo, exporter, notifier)
	
	// Define yesterday's time range
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	startOfYesterday := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())
	endOfYesterday := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 999999999, yesterday.Location())
	
	// Create test entries
	entry1 := entities.OvertimeEntry{
		TicketURL: "http://jira.com/ticket1",
		Minutes:   90,
		Date:      yesterday,
	}
	
	entry2 := entities.OvertimeEntry{
		TicketURL: "http://jira.com/ticket2",
		Minutes:   60,
		Date:      yesterday,
	}
	
	// Add test entries to repository
	repo.AddTestEntry(entry1, startOfYesterday, endOfYesterday)
	repo.AddTestEntry(entry2, startOfYesterday, endOfYesterday)
	
	// Execute the use case
	ctx := context.Background()
	err := uc.ProcessYesterdayOvertime(ctx)
	
	// Check results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Check repository
	currentMonth := now.Format("Jan-2006")
	report, err := repo.GetMergedReport(ctx, currentMonth)
	if err != nil {
		t.Fatalf("Error getting merged report: %v", err)
	}
	
	if len(report.Entries) != 2 {
		t.Errorf("Expected 2 entries in report, got %d", len(report.Entries))
	}
	
	if report.TotalTime != 150 {
		t.Errorf("Expected total time 150, got %d", report.TotalTime)
	}
}

func TestProcessYesterdayOvertimeError(t *testing.T) {
	// Create mocks
	repo := mocks.NewMockOvertimeRepository()
	exporter := mocks.NewMockReportExporter()
	notifier := mocks.NewMockNotificationService()
	
	// Set up error
	testError := errors.New("test error")
	repo.GetPeriodError = testError
	
	// Create use case
	uc := usecases.NewOvertimeUseCase(repo, exporter, notifier)
	
	// Execute the use case
	ctx := context.Background()
	err := uc.ProcessYesterdayOvertime(ctx)
	
	// Check error was propagated
	if err == nil {
		t.Error("Expected error, got nil")
	} else if !errors.Is(err, testError) {
		t.Errorf("Expected error %v, got %v", testError, err)
	}
}

func TestGenerateMonthlyReport(t *testing.T) {
	// Create mocks
	repo := mocks.NewMockOvertimeRepository()
	exporter := mocks.NewMockReportExporter()
	notifier := mocks.NewMockNotificationService()
	
	// Create use case
	uc := usecases.NewOvertimeUseCase(repo, exporter, notifier)
	
	// Set up test data
	now := time.Now()
	prevMonth := now.AddDate(0, -1, 0)
	monthPeriod := prevMonth.Format("Jan-2006")
	
	report := entities.NewOvertimeReport(monthPeriod)
	report.AddEntry("http://jira.com/ticket1", 120)
	report.AddEntry("http://jira.com/ticket2", 45)
	
	// Add test report to repository
	repo.AddTestReport(report)
	
	// Execute the use case
	ctx := context.Background()
	err := uc.GenerateMonthlyReport(ctx)
	
	// Check results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Check exporter was called
	if exporter.ExportToExcelCalls != 1 {
		t.Errorf("Expected ExportToExcel to be called once, got %d", exporter.ExportToExcelCalls)
	}
	
	// Check notification service was called
	if notifier.SendEmailCalls != 1 {
		t.Errorf("Expected SendEmail to be called once, got %d", notifier.SendEmailCalls)
	}
	
	// Check notification service received correct parameters
	if notifier.LastReportSent != report {
		t.Error("Expected correct report to be sent to notification service")
	}
	
	if notifier.LastAttachmentPath != exporter.ExcelFilePath {
		t.Errorf("Expected attachment path %s, got %s", exporter.ExcelFilePath, notifier.LastAttachmentPath)
	}
}

func TestTestMonthlyReport(t *testing.T) {
	// Create mocks
	repo := mocks.NewMockOvertimeRepository()
	exporter := mocks.NewMockReportExporter()
	notifier := mocks.NewMockNotificationService()
	
	// Create use case
	uc := usecases.NewOvertimeUseCase(repo, exporter, notifier)
	
	// Set up test data
	now := time.Now()
	currentMonth := now.Format("Jan-2006")
	
	report := entities.NewOvertimeReport(currentMonth)
	report.AddEntry("http://jira.com/ticket1", 60)
	report.AddEntry("http://jira.com/ticket2", 30)
	
	// Add test report to repository
	repo.AddTestReport(report)
	
	// Execute the use case
	ctx := context.Background()
	err := uc.TestMonthlyReport(ctx)
	
	// Check results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Check exporter was called
	if exporter.ExportToExcelCalls != 1 {
		t.Errorf("Expected ExportToExcel to be called once, got %d", exporter.ExportToExcelCalls)
	}
	
	// Check notification service was called
	if notifier.SendEmailCalls != 1 {
		t.Errorf("Expected SendEmail to be called once, got %d", notifier.SendEmailCalls)
	}
	
	// Check notification service received correct parameters
	if notifier.LastReportSent != report {
		t.Error("Expected correct report to be sent to notification service")
	}
	
	if notifier.LastAttachmentPath != exporter.ExcelFilePath {
		t.Errorf("Expected attachment path %s, got %s", exporter.ExcelFilePath, notifier.LastAttachmentPath)
	}
}