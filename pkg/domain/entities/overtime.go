package entities

import "time"

// OvertimeEntry represents a single overtime record
type OvertimeEntry struct {
	TicketURL string
	Minutes   int
	Date      time.Time
}

// OvertimeReport represents a collection of overtime entries for a reporting period
type OvertimeReport struct {
	Entries    []OvertimeEntry
	Period     string
	TotalTime  int
	ReportDate time.Time
}

// CalculateTotalMinutes computes the total minutes from all entries
func (r *OvertimeReport) CalculateTotalMinutes() int {
	total := 0
	for _, entry := range r.Entries {
		total += entry.Minutes
	}
	r.TotalTime = total
	return total
}

// NewOvertimeReport creates a new overtime report for the given period
func NewOvertimeReport(period string) *OvertimeReport {
	return &OvertimeReport{
		Entries:    []OvertimeEntry{},
		Period:     period,
		ReportDate: time.Now(),
	}
}

// AddEntry adds a new overtime entry to the report
func (r *OvertimeReport) AddEntry(ticketURL string, minutes int) {
	entry := OvertimeEntry{
		TicketURL: ticketURL,
		Minutes:   minutes,
		Date:      time.Now(),
	}
	r.Entries = append(r.Entries, entry)
	r.CalculateTotalMinutes()
}