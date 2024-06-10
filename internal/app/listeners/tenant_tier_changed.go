package listeners

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/events"
	"github.com/Marcellinom/tenant-management-saas/provider/event"
	"time"
)

type TenantTierChangedListener struct {
}

func NewTenantTierChangedListener() TenantTierChangedListener {
	return TenantTierChangedListener{}
}

func (receiver TenantTierChangedListener) Name() string {
	return fmt.Sprintf("%T", receiver)
}

func (receiver TenantTierChangedListener) Handle(ctx context.Context, event event.Event) error {
	select {
	case <-ctx.Value("stop-signal").(chan bool):
		return ctx.Err()
	default:
		time.Sleep(5 * time.Second)
		var payload events.TenantTierChanged
		json_data, err := event.JSON()
		if err != nil {
			return fmt.Errorf("gagal menencode json pada event listener: %w", err)
		}
		err = json.Unmarshal(json_data, &payload)
		if err != nil {
			return fmt.Errorf("gagal mendecode json pada event listener: %w", err)
		}
		fmt.Println("success memroses event!", payload)
		return nil
	}
}

func (receiver TenantTierChangedListener) MaxRetries() int {
	return 3
}
