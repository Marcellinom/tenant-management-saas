package listeners

import (
	"context"
	"fmt"
	"github.com/Marcellinom/tenant-management-saas/provider/event"
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
	fmt.Println(event.JSON())
	return nil
}

func (receiver TenantTierChangedListener) MaxRetries() int {
	return 3
}
