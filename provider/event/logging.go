package event

import (
	"github.com/Marcellinom/tenant-management-saas/provider"
	"time"
)

func MarkAsFailed(app *provider.Application, event_name, listener_name, message string, metadata []byte, max_retries int, fail_type ...string) {
	payload := make(map[string]any)
	payload = map[string]any{
		"event_name":    event_name,
		"listener_name": listener_name,
		"metadata":      metadata,
		"max_retries":   max_retries,
		"created_at":    time.Now(),
		"message":       message,
	}
	if len(fail_type) > 0 {
		payload["type"] = fail_type[0]
	}
	app.DefaultDatabase().Table("failed_jobs").Create(payload)
}
