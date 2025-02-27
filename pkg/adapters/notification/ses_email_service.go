package notification

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/MateSousa/overtime-script/pkg/domain/entities"
	"github.com/MateSousa/overtime-script/pkg/domain/repositories"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// SESEmailService implements the NotificationService interface using AWS SES
type SESEmailService struct {
	senderEmail string
	recipient   string
	region      string
}

// NewSESEmailService creates a new AWS SES email service
func NewSESEmailService(senderEmail, recipient, region string) repositories.NotificationService {
	return &SESEmailService{
		senderEmail: senderEmail,
		recipient:   recipient,
		region:      region,
	}
}

// SendReportByEmail sends an overtime report via email with an attachment
func (s *SESEmailService) SendReportByEmail(ctx context.Context, report *entities.OvertimeReport, attachmentPath string) error {
	// Build the email body with a simple message
	emailBody := fmt.Sprintf(`Caros,

Espero que estejam bem!

Segue em anexo as horas extra do mês de %s.

Atenciosamente,`, report.Period)

	// Read the file content
	fileContent, err := os.ReadFile(attachmentPath)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(s.region),
	})
	if err != nil {
		return fmt.Errorf("error creating AWS session: %w", err)
	}

	// Create multipart message boundary
	boundary := "==Multipart_Boundary_x" + time.Now().Format("20060102150405") + "x"

	// Determine content type based on file extension
	fileExt := filepath.Ext(attachmentPath)
	contentType := "application/octet-stream" // Default content type
	
	if fileExt == ".csv" {
		contentType = "text/csv"
	} else if fileExt == ".xlsx" {
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	}

	// Create the raw message
	rawMessage := fmt.Sprintf("From: %s\n", s.senderEmail) +
		fmt.Sprintf("To: %s\n", s.recipient) +
		"Subject: Darede - Relatório Mensal de Horas Extras\n" +
		"MIME-Version: 1.0\n" +
		fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\n", boundary) +
		"\n" +
		fmt.Sprintf("--%s\n", boundary) +
		"Content-Type: text/plain; charset=UTF-8\n" +
		"Content-Transfer-Encoding: 7bit\n" +
		"\n" +
		emailBody + "\n" +
		"\n" +
		fmt.Sprintf("--%s\n", boundary) +
		fmt.Sprintf("Content-Type: %s; charset=UTF-8\n", contentType) +
		"Content-Transfer-Encoding: base64\n" +
		fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\n", filepath.Base(attachmentPath)) +
		"\n" +
		base64.StdEncoding.EncodeToString(fileContent) + "\n" +
		"\n" +
		fmt.Sprintf("--%s--", boundary)

	// Create an SES client
	svc := ses.New(sess)
	input := &ses.SendRawEmailInput{
		Destinations: []*string{
			aws.String(s.recipient),
		},
		RawMessage: &ses.RawMessage{
			Data: []byte(rawMessage),
		},
		Source: aws.String(s.senderEmail),
	}

	// Send the email
	_, err = svc.SendRawEmail(input)
	if err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	return nil
}