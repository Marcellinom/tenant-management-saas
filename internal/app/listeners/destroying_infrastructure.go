package listeners

import (
	"context"
	"fmt"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/repositories"
	"github.com/Marcellinom/tenant-management-saas/provider/event"
)

type DestroyingInfrastructureListener struct {
	infra_service repositories.InfrastructureRepositoryInterface
}

func NewDestroyingInfrastructureListener(infra_service repositories.InfrastructureRepositoryInterface) *DestroyingInfrastructureListener {
	return &DestroyingInfrastructureListener{infra_service: infra_service}
}

func (l DestroyingInfrastructureListener) Handle(ctx context.Context, event event.Event) error {
	// TODO: implement this
	return nil
}

func (l DestroyingInfrastructureListener) MaxRetries() int {
	return 3
}

func (l DestroyingInfrastructureListener) Name() string {
	return fmt.Sprintf("%T", l)
}
