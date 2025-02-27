package repositories

import (
	"context"
	"time"

	"github.com/MateSousa/overtime-script/pkg/domain/entities"
)

// OvertimeRepository defines the interface for accessing overtime data
type OvertimeRepository interface {
	// GetOvertimeEntriesForPeriod fetches overtime entries for a specific time period
	GetOvertimeEntriesForPeriod(ctx context.Context, start, end time.Time) ([]entities.OvertimeEntry, error)
	
	// SaveOvertimeReport persists an overtime report
	SaveOvertimeReport(ctx context.Context, report *entities.OvertimeReport) error
	
	// GetMergedReport retrieves the merged overtime report for a specific month
	GetMergedReport(ctx context.Context, month string) (*entities.OvertimeReport, error)
	
	// MergeOvertimeEntries combines multiple overtime entries into a single report
	MergeOvertimeEntries(ctx context.Context, entries []entities.OvertimeEntry, period string) (*entities.OvertimeReport, error)
}