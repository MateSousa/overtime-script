package mocks

import (
	"context"
	"fmt"
	"time"

	"github.com/MateSousa/overtime-script/pkg/domain/entities"
)

// MockOvertimeRepository is a mock implementation of the OvertimeRepository interface
type MockOvertimeRepository struct {
	entries            map[string][]entities.OvertimeEntry
	reports            map[string]*entities.OvertimeReport
	ErrorToReturn      error
	GetPeriodError     error
	SaveReportError    error
	GetMergedError     error
	MergeEntriesError  error
}

// NewMockOvertimeRepository creates a new mock repository
func NewMockOvertimeRepository() *MockOvertimeRepository {
	return &MockOvertimeRepository{
		entries: make(map[string][]entities.OvertimeEntry),
		reports: make(map[string]*entities.OvertimeReport),
	}
}

// GetOvertimeEntriesForPeriod fetches overtime entries for a specific time period
func (m *MockOvertimeRepository) GetOvertimeEntriesForPeriod(ctx context.Context, start, end time.Time) ([]entities.OvertimeEntry, error) {
	if m.GetPeriodError != nil {
		return nil, m.GetPeriodError
	}

	// Generate key from start and end dates
	key := fmt.Sprintf("%s-%s", start.Format("2006-01-02"), end.Format("2006-01-02"))
	
	// Return entries for this period, or empty slice if none
	entries, ok := m.entries[key]
	if !ok {
		return []entities.OvertimeEntry{}, nil
	}
	
	return entries, nil
}

// SaveOvertimeReport persists an overtime report
func (m *MockOvertimeRepository) SaveOvertimeReport(ctx context.Context, report *entities.OvertimeReport) error {
	if m.SaveReportError != nil {
		return m.SaveReportError
	}
	
	// Store the report
	m.reports[report.Period] = report
	return nil
}

// GetMergedReport retrieves the merged overtime report for a specific month
func (m *MockOvertimeRepository) GetMergedReport(ctx context.Context, month string) (*entities.OvertimeReport, error) {
	if m.GetMergedError != nil {
		return nil, m.GetMergedError
	}
	
	report, ok := m.reports[month]
	if !ok {
		return nil, fmt.Errorf("report not found for period: %s", month)
	}
	
	return report, nil
}

// MergeOvertimeEntries combines multiple overtime entries into a single report
func (m *MockOvertimeRepository) MergeOvertimeEntries(ctx context.Context, entries []entities.OvertimeEntry, period string) (*entities.OvertimeReport, error) {
	if m.MergeEntriesError != nil {
		return nil, m.MergeEntriesError
	}
	
	// Try to get existing report
	existingReport, ok := m.reports[period]
	if !ok {
		// Create new report if not found
		existingReport = entities.NewOvertimeReport(period)
	}
	
	// Add entries to the report
	for _, entry := range entries {
		existingReport.AddEntry(entry.TicketURL, entry.Minutes)
	}
	
	// Save the updated report
	m.reports[period] = existingReport
	
	return existingReport, nil
}

// AddTestEntry adds a test entry to the repository
func (m *MockOvertimeRepository) AddTestEntry(entry entities.OvertimeEntry, startDate, endDate time.Time) {
	key := fmt.Sprintf("%s-%s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	m.entries[key] = append(m.entries[key], entry)
}

// AddTestReport adds a test report to the repository
func (m *MockOvertimeRepository) AddTestReport(report *entities.OvertimeReport) {
	m.reports[report.Period] = report
}