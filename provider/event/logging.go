package event

import (
	"github.com/Marcellinom/tenant-management-saas/provider"
	"time"
)

func MarkAsFailed(app *provider.Application, event_name, listener_name, message string, metadata []byte, max_retries int) {
	app.DefaultDatabase().Table("failed_jobs").Create(map[string]any{
		"event_name":    event_name,
		"listener_name": listener_name,
		"metadata":      metadata,
		"max_retries":   max_retries,
		"created_at":    time.Now(),
		"message":       message,
	})
}