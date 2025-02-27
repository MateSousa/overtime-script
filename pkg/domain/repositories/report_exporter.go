package repositories

import (
	"context"

	"github.com/MateSousa/overtime-script/pkg/domain/entities"
)

// ReportExporter defines the interface for exporting reports to different formats
type ReportExporter interface {
	// ExportToExcel exports the report to an Excel file
	ExportToExcel(ctx context.Context, report *entities.OvertimeReport) (string, error)
	
	// ExportToCSV exports the report to a CSV file (for backward compatibility)
	ExportToCSV(ctx context.Context, report *entities.OvertimeReport) (string, error)
}