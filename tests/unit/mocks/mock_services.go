package mocks

import (
	"context"

	"github.com/MateSousa/overtime-script/pkg/domain/entities"
)

// MockReportExporter is a mock implementation of the ReportExporter interface
type MockReportExporter struct {
	ExcelFilePath      string
	CSVFilePath        string
	ExportToExcelError error
	ExportToCSVError   error
	ExportToExcelCalls int
	ExportToCSVCalls   int
}

// NewMockReportExporter creates a new mock report exporter
func NewMockReportExporter() *MockReportExporter {
	return &MockReportExporter{
		ExcelFilePath: "test-report.xlsx",
		CSVFilePath:   "test-report.csv",
	}
}

// ExportToExcel exports the report to an Excel file
func (m *MockReportExporter) ExportToExcel(ctx context.Context, report *entities.OvertimeReport) (string, error) {
	m.ExportToExcelCalls++
	if m.ExportToExcelError != nil {
		return "", m.ExportToExcelError
	}
	return m.ExcelFilePath, nil
}

// ExportToCSV exports the report to a CSV file
func (m *MockReportExporter) ExportToCSV(ctx context.Context, report *entities.OvertimeReport) (string, error) {
	m.ExportToCSVCalls++
	if m.ExportToCSVError != nil {
		return "", m.ExportToCSVError
	}
	return m.CSVFilePath, nil
}

// MockNotificationService is a mock implementation of the NotificationService interface
type MockNotificationService struct {
	SendEmailError   error
	SendEmailCalls   int
	LastReportSent   *entities.OvertimeReport
	LastAttachmentPath string
}

// NewMockNotificationService creates a new mock notification service
func NewMockNotificationService() *MockNotificationService {
	return &MockNotificationService{}
}

// SendReportByEmail sends an overtime report via email with an attachment
func (m *MockNotificationService) SendReportByEmail(ctx context.Context, report *entities.OvertimeReport, attachmentPath string) error {
	m.SendEmailCalls++
	m.LastReportSent = report
	m.LastAttachmentPath = attachmentPath
	return m.SendEmailError
}