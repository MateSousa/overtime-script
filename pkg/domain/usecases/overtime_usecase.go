package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/MateSousa/overtime-script/pkg/domain/repositories"
)

// OvertimeUseCase defines the overtime business logic
type OvertimeUseCase struct {
	repository         repositories.OvertimeRepository
	reportExporter     repositories.ReportExporter
	notificationService repositories.NotificationService
}

// NewOvertimeUseCase creates a new overtime use case instance
func NewOvertimeUseCase(
	repo repositories.OvertimeRepository,
	exporter repositories.ReportExporter,
	notifier repositories.NotificationService,
) *OvertimeUseCase {
	return &OvertimeUseCase{
		repository:         repo,
		reportExporter:     exporter,
		notificationService: notifier,
	}
}

// ProcessYesterdayOvertime collects and processes overtime entries from yesterday
func (uc *OvertimeUseCase) ProcessYesterdayOvertime(ctx context.Context) error {
	// Define yesterday's time range
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	startOfYesterday := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())
	endOfYesterday := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 999999999, yesterday.Location())
	
	// Get yesterday's overtime entries
	entries, err := uc.repository.GetOvertimeEntriesForPeriod(ctx, startOfYesterday, endOfYesterday)
	if err != nil {
		return fmt.Errorf("error getting yesterday's overtime entries: %w", err)
	}
	
	// Create and save the monthly report
	currentMonthPeriod := now.Format("Jan-2006")
	report, err := uc.repository.MergeOvertimeEntries(ctx, entries, currentMonthPeriod)
	if err != nil {
		return fmt.Errorf("error merging overtime entries: %w", err)
	}
	
	// Save the updated report
	if err := uc.repository.SaveOvertimeReport(ctx, report); err != nil {
		return fmt.Errorf("error saving overtime report: %w", err)
	}
	
	return nil
}

// GenerateMonthlyReport generates the report for the previous month and sends it via email
func (uc *OvertimeUseCase) GenerateMonthlyReport(ctx context.Context) error {
	// Calculate previous month
	now := time.Now()
	prevMonth := now.AddDate(0, -1, 0)
	monthPeriod := prevMonth.Format("Jan-2006")
	
	// Get the merged report for the previous month
	report, err := uc.repository.GetMergedReport(ctx, monthPeriod)
	if err != nil {
		return fmt.Errorf("error getting merged report for %s: %w", monthPeriod, err)
	}
	
	// Export the report to Excel
	excelFilePath, err := uc.reportExporter.ExportToExcel(ctx, report)
	if err != nil {
		return fmt.Errorf("error exporting report to Excel: %w", err)
	}
	
	// Send the report via email
	if err := uc.notificationService.SendReportByEmail(ctx, report, excelFilePath); err != nil {
		return fmt.Errorf("error sending report email: %w", err)
	}
	
	return nil
}

// TestMonthlyReport generates a test report for the current month and sends it via email
func (uc *OvertimeUseCase) TestMonthlyReport(ctx context.Context) error {
	// Get current month
	now := time.Now()
	monthPeriod := now.Format("Jan-2006")
	
	// Get the merged report for the current month
	report, err := uc.repository.GetMergedReport(ctx, monthPeriod)
	if err != nil {
		return fmt.Errorf("error getting merged report for %s: %w", monthPeriod, err)
	}
	
	// Export the report to Excel
	excelFilePath, err := uc.reportExporter.ExportToExcel(ctx, report)
	if err != nil {
		return fmt.Errorf("error exporting report to Excel: %w", err)
	}
	
	// Send the report via email
	if err := uc.notificationService.SendReportByEmail(ctx, report, excelFilePath); err != nil {
		return fmt.Errorf("error sending report email: %w", err)
	}
	
	return nil
}