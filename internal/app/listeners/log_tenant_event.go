package listeners

import (
	"context"
	"fmt"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/repositories"
	"github.com/Marcellinom/tenant-management-saas/provider/event"
)

type LogTenantEvent struct {
	tenant_repo repositories.TenantRepositoryInterface
}

func NewLogTenantEvent(tenant_repo repositories.TenantRepositoryInterface) *LogTenantEvent {
	return &LogTenantEvent{tenant_repo: tenant_repo}
}

func (l LogTenantEvent) Handle(ctx context.Context, event event.Event) error {
	// TODO: implement this
	return nil
}

func (l LogTenantEvent) MaxRetries() int {
	return 3
}

func (l LogTenantEvent) Name() string {
	return fmt.Sprintf("%T", l)
}
