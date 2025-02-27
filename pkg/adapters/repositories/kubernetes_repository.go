package repositories

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/MateSousa/overtime-script/pkg/domain/entities"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// KubernetesOvertimeRepository implements the OvertimeRepository interface using Kubernetes ConfigMaps
type KubernetesOvertimeRepository struct {
	client    *kubernetes.Clientset
	namespace string
}

// NewKubernetesOvertimeRepository creates a new Kubernetes repository instance
func NewKubernetesOvertimeRepository(client *kubernetes.Clientset, namespace string) *KubernetesOvertimeRepository {
	return &KubernetesOvertimeRepository{
		client:    client,
		namespace: namespace,
	}
}

// GetOvertimeEntriesForPeriod fetches overtime entries for a specific time period from ConfigMaps
func (r *KubernetesOvertimeRepository) GetOvertimeEntriesForPeriod(ctx context.Context, start, end time.Time) ([]entities.OvertimeEntry, error) {
	// List ConfigMaps with label selector "app=overtime"
	list, err := r.client.CoreV1().ConfigMaps(r.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: "app=overtime",
	})
	if err != nil {
		return nil, fmt.Errorf("error listing ConfigMaps: %w", err)
	}

	var entries []entities.OvertimeEntry
	for _, cm := range list.Items {
		// Check if the ConfigMap was created in the requested time period
		if (cm.CreationTimestamp.Time.After(start) || cm.CreationTimestamp.Time.Equal(start)) && 
		   (cm.CreationTimestamp.Time.Before(end) || cm.CreationTimestamp.Time.Equal(end)) {
			
			tickets, ticketsOk := cm.Data["ticket_url"]
			minutes, minutesOk := cm.Data["minutes"]
			
			if !ticketsOk || !minutesOk {
				continue
			}
			
			ticketList := strings.Split(tickets, "\n")
			minutesList := strings.Split(minutes, "\n")
			
			// Use the smaller length if counts differ
			count := len(ticketList)
			if len(minutesList) < count {
				count = len(minutesList)
			}
			
			for i := 0; i < count; i++ {
				ticket := strings.TrimSpace(ticketList[i])
				minuteStr := strings.TrimSpace(minutesList[i])
				minuteVal, err := strconv.Atoi(minuteStr)
				if err != nil {
					minuteVal = 0
				}
				
				entry := entities.OvertimeEntry{
					TicketURL: ticket,
					Minutes:   minuteVal,
					Date:      cm.CreationTimestamp.Time,
				}
				entries = append(entries, entry)
			}
		}
	}

	return entries, nil
}

// SaveOvertimeReport saves an overtime report as a ConfigMap
func (r *KubernetesOvertimeRepository) SaveOvertimeReport(ctx context.Context, report *entities.OvertimeReport) error {
	// Create data map for ConfigMap
	data := make(map[string]string)
	
	var ticketURLs, minutes []string
	for _, entry := range report.Entries {
		ticketURLs = append(ticketURLs, entry.TicketURL)
		minutes = append(minutes, strconv.Itoa(entry.Minutes))
	}
	
	data["ticket_url"] = strings.Join(ticketURLs, "\n")
	data["minutes"] = strings.Join(minutes, "\n")
	
	// ConfigMap name (lowercase for RFC1123 compliance)
	cmName := strings.ToLower(report.Period + "-overtime-merged")
	
	// Try to get existing ConfigMap
	cmInterface := r.client.CoreV1().ConfigMaps(r.namespace)
	existing, err := cmInterface.Get(ctx, cmName, metav1.GetOptions{})
	
	if err != nil {
		// Not found - create a new one
		newCM := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name: cmName,
			},
			Data: data,
		}
		_, err = cmInterface.Create(ctx, newCM, metav1.CreateOptions{})
		return err
	}
	
	// Update existing ConfigMap
	existing.Data = data
	_, err = cmInterface.Update(ctx, existing, metav1.UpdateOptions{})
	return err
}

// GetMergedReport retrieves the merged overtime report for a specific month
func (r *KubernetesOvertimeRepository) GetMergedReport(ctx context.Context, month string) (*entities.OvertimeReport, error) {
	// ConfigMap name (lowercase for RFC1123 compliance)
	cmName := strings.ToLower(month + "-overtime-merged")
	
	// Get the ConfigMap
	cm, err := r.client.CoreV1().ConfigMaps(r.namespace).Get(ctx, cmName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("error getting merged ConfigMap %s: %w", cmName, err)
	}
	
	report := entities.NewOvertimeReport(month)
	
	tickets, ticketsOk := cm.Data["ticket_url"]
	minutes, minutesOk := cm.Data["minutes"]
	
	if !ticketsOk || !minutesOk {
		return report, nil // Return empty report if no data
	}
	
	ticketList := strings.Split(tickets, "\n")
	minutesList := strings.Split(minutes, "\n")
	
	// Use the smaller length if counts differ
	count := len(ticketList)
	if len(minutesList) < count {
		count = len(minutesList)
	}
	
	for i := 0; i < count; i++ {
		ticket := strings.TrimSpace(ticketList[i])
		minuteStr := strings.TrimSpace(minutesList[i])
		minuteVal, err := strconv.Atoi(minuteStr)
		if err != nil {
			minuteVal = 0
		}
		
		report.AddEntry(ticket, minuteVal)
	}
	
	return report, nil
}

// MergeOvertimeEntries combines multiple overtime entries into a single report
func (r *KubernetesOvertimeRepository) MergeOvertimeEntries(ctx context.Context, entries []entities.OvertimeEntry, period string) (*entities.OvertimeReport, error) {
	// First try to get existing report for the period
	existingReport, err := r.GetMergedReport(ctx, period)
	if err != nil {
		// If not found, create a new report
		existingReport = entities.NewOvertimeReport(period)
	}
	
	// Add all new entries to the existing report
	for _, entry := range entries {
		existingReport.AddEntry(entry.TicketURL, entry.Minutes)
	}
	
	return existingReport, nil
}