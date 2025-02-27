package repositories

import (
	"context"

	"github.com/MateSousa/overtime-script/pkg/domain/entities"
)

// NotificationService defines the interface for sending notifications
type NotificationService interface {
	// SendReportByEmail sends an overtime report via email with an attachment
	SendReportByEmail(ctx context.Context, report *entities.OvertimeReport, attachmentPath string) error
}