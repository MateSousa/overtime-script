package unit

import (
	"testing"
	"time"

	"github.com/MateSousa/overtime-script/pkg/domain/entities"
)

func TestNewOvertimeReport(t *testing.T) {
	period := "Jan-2023"
	report := entities.NewOvertimeReport(period)

	if report.Period != period {
		t.Errorf("Expected period %s, got %s", period, report.Period)
	}

	if report.TotalTime != 0 {
		t.Errorf("Expected total time 0, got %d", report.TotalTime)
	}

	if len(report.Entries) != 0 {
		t.Errorf("Expected empty entries, got %d entries", len(report.Entries))
	}
}

func TestAddEntry(t *testing.T) {
	report := entities.NewOvertimeReport("Jan-2023")

	// Add first entry
	report.AddEntry("http://ticket1.com", 120)

	if len(report.Entries) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(report.Entries))
	}

	if report.TotalTime != 120 {
		t.Errorf("Expected total time 120, got %d", report.TotalTime)
	}

	// Add second entry
	report.AddEntry("http://ticket2.com", 60)

	if len(report.Entries) != 2 {
		t.Fatalf("Expected 2 entries, got %d", len(report.Entries))
	}

	if report.TotalTime != 180 {
		t.Errorf("Expected total time 180, got %d", report.TotalTime)
	}

	// Check entries
	if report.Entries[0].TicketURL != "http://ticket1.com" {
		t.Errorf("Expected first ticket URL http://ticket1.com, got %s", report.Entries[0].TicketURL)
	}

	if report.Entries[0].Minutes != 120 {
		t.Errorf("Expected first ticket minutes 120, got %d", report.Entries[0].Minutes)
	}

	if report.Entries[1].TicketURL != "http://ticket2.com" {
		t.Errorf("Expected second ticket URL http://ticket2.com, got %s", report.Entries[1].TicketURL)
	}

	if report.Entries[1].Minutes != 60 {
		t.Errorf("Expected second ticket minutes 60, got %d", report.Entries[1].Minutes)
	}
}

func TestCalculateTotalMinutes(t *testing.T) {
	report := entities.NewOvertimeReport("Jan-2023")

	// Empty report should have 0 total
	total := report.CalculateTotalMinutes()
	if total != 0 {
		t.Errorf("Expected total minutes 0 for empty report, got %d", total)
	}

	// Add entries
	report.Entries = append(report.Entries, entities.OvertimeEntry{
		TicketURL: "http://ticket1.com",
		Minutes:   45,
		Date:      time.Now(),
	})

	report.Entries = append(report.Entries, entities.OvertimeEntry{
		TicketURL: "http://ticket2.com",
		Minutes:   30,
		Date:      time.Now(),
	})

	report.Entries = append(report.Entries, entities.OvertimeEntry{
		TicketURL: "http://ticket3.com",
		Minutes:   75,
		Date:      time.Now(),
	})

	// Recalculate total
	total = report.CalculateTotalMinutes()
	if total != 150 {
		t.Errorf("Expected total minutes 150, got %d", total)
	}

	// Check if total was saved in the struct
	if report.TotalTime != 150 {
		t.Errorf("Expected TotalTime field to be 150, got %d", report.TotalTime)
	}
}