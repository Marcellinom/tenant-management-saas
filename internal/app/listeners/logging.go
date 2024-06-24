package listeners

import (
	"context"
	"fmt"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/repositories"
	"github.com/Marcellinom/tenant-management-saas/provider/event"
)

type Logging struct {
	tenant_repo repositories.TenantRepositoryInterface
}

func LogTenantEvent(tenant_repo repositories.TenantRepositoryInterface) *Logging {
	return &Logging{tenant_repo: tenant_repo}
}

func (l Logging) Handle(ctx context.Context, event event.Event) error {
	// TODO: implement this
	return nil
}

func (l Logging) MaxRetries() int {
	return 3
}

func (l Logging) Name() string {
	return fmt.Sprintf("%T", l)
}
