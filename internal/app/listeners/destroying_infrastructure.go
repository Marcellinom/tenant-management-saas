package listeners

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/events"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/repositories"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/services"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/vo"
	"github.com/Marcellinom/tenant-management-saas/provider/event"
)

type DestroyingInfrastructureListener struct {
	infra_repo       repositories.InfrastructureRepositoryInterface
	deployer_service services.DeployerServiceInterface
}

func NewDestroyingInfrastructureListener(
	infra_repo repositories.InfrastructureRepositoryInterface,
	deployer_service services.DeployerServiceInterface,
) *DestroyingInfrastructureListener {
	return &DestroyingInfrastructureListener{
		infra_repo:       infra_repo,
		deployer_service: deployer_service,
	}
}

func (l DestroyingInfrastructureListener) Handle(ctx context.Context, event event.Event) error {
	var payload events.InfrastructureDestroyed
	json_data, err := event.JSON()
	if err != nil {
		return fmt.Errorf("gagal menencode json pada event listener: %w", err)
	}
	err = json.Unmarshal(json_data, &payload)
	if err != nil {
		return fmt.Errorf("gagal mendecode json pada event listener: %w", err)
	}

	infra_id, err := vo.NewInfrastructureId(payload.InfrastructureId)
	if err != nil {
		return fmt.Errorf("terjadi kesalahan dalam parsing id infrastructure (wants uuid, get %s)", payload.InfrastructureId)
	}
	infra, err := l.infra_repo.Find(infra_id)

	return l.deployer_service.DecommissionInfrastructure(ctx, infra)
}

func (l DestroyingInfrastructureListener) MaxRetries() int {
	return 3
}

func (l DestroyingInfrastructureListener) Name() string {
	return fmt.Sprintf("%T", l)
}
