package exporters

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/MateSousa/overtime-script/pkg/domain/entities"
	"github.com/MateSousa/overtime-script/pkg/domain/repositories"
	"github.com/xuri/excelize/v2"
)

// ExcelReportExporter implements the ReportExporter interface for Excel files
type ExcelReportExporter struct{}

// NewExcelReportExporter creates a new Excel report exporter
func NewExcelReportExporter() repositories.ReportExporter {
	return &ExcelReportExporter{}
}

// ExportToExcel exports the report to an Excel file
func (e *ExcelReportExporter) ExportToExcel(ctx context.Context, report *entities.OvertimeReport) (string, error) {
	// Create a new Excel file
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("Error closing Excel file:", err)
		}
	}()

	// Get the default sheet (usually "Sheet1")
	sheetName := f.GetSheetName(0)
	
	// Set column widths
	f.SetColWidth(sheetName, "A", "A", 50) // Ticket column width
	f.SetColWidth(sheetName, "B", "B", 15) // Minutes column width
	
	// Create header style - Blue background with white text, bold, centered
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Color:  "#FFFFFF",
			Size:   12,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
		},
	})
	if err != nil {
		return "", fmt.Errorf("error creating header style: %w", err)
	}
	
	// Create data style - Light borders
	dataStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
		},
	})
	if err != nil {
		return "", fmt.Errorf("error creating data style: %w", err)
	}
	
	// Create total row style - Bold with gray background
	totalStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#D9D9D9"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
		},
	})
	if err != nil {
		return "", fmt.Errorf("error creating total style: %w", err)
	}
	
	// Write headers
	f.SetCellValue(sheetName, "A1", "TICKET")
	f.SetCellValue(sheetName, "B1", "MINUTOS")
	
	// Apply header style
	f.SetCellStyle(sheetName, "A1", "B1", headerStyle)
	
	// Write data rows
	for i, entry := range report.Entries {
		rowNum := i + 2 // Start from row 2 (after headers)
		
		// Set values
		cellA := fmt.Sprintf("A%d", rowNum)
		cellB := fmt.Sprintf("B%d", rowNum)
		f.SetCellValue(sheetName, cellA, entry.TicketURL)
		f.SetCellValue(sheetName, cellB, entry.Minutes)
		
		// Apply data style
		f.SetCellStyle(sheetName, cellA, cellB, dataStyle)
	}
	
	// Write total row
	totalRow := len(report.Entries) + 2
	totalCellA := fmt.Sprintf("A%d", totalRow)
	totalCellB := fmt.Sprintf("B%d", totalRow)
	f.SetCellValue(sheetName, totalCellA, "TOTAL")
	f.SetCellValue(sheetName, totalCellB, report.TotalTime)
	
	// Apply total style
	f.SetCellStyle(sheetName, totalCellA, totalCellB, totalStyle)
	
	// Create file with the report period
	filename := fmt.Sprintf("overtime_%s.xlsx", time.Now().Format("2006-01-02"))
	if err := f.SaveAs(filename); err != nil {
		return "", fmt.Errorf("error saving Excel file: %w", err)
	}

	return filename, nil
}

// ExportToCSV exports the report to a CSV file (for backward compatibility)
func (e *ExcelReportExporter) ExportToCSV(ctx context.Context, report *entities.OvertimeReport) (string, error) {
	// Create CSV file with the current date
	filename := fmt.Sprintf("overtime_%s.csv", time.Now().Format("2006-01-02"))
	file, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("error creating CSV file: %w", err)
	}
	defer file.Close()

	// Create Excel-compatible CSV with semicolons as separators
	writer := csv.NewWriter(file)
	writer.Comma = ';' // Use semicolon as separator for Excel compatibility
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"TICKET", "MINUTOS"}); err != nil {
		return "", fmt.Errorf("error writing CSV header: %w", err)
	}

	// Write data rows
	for _, entry := range report.Entries {
		if err := writer.Write([]string{entry.TicketURL, strconv.Itoa(entry.Minutes)}); err != nil {
			return "", fmt.Errorf("error writing CSV row: %w", err)
		}
	}

	// Write total row
	if err := writer.Write([]string{"TOTAL", strconv.Itoa(report.TotalTime)}); err != nil {
		return "", fmt.Errorf("error writing CSV total row: %w", err)
	}

	return filename, nil
}