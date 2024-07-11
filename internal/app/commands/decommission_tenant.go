package commands

import (
	"fmt"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/repositories"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/vo"
	"github.com/Marcellinom/tenant-management-saas/provider/errors"
)

type DecommissionTenantCommand struct {
	tenant_repo repositories.TenantRepositoryInterface
}

func NewDecommissionTenantCommand(tenant_repo repositories.TenantRepositoryInterface) *DecommissionTenantCommand {
	return &DecommissionTenantCommand{tenant_repo: tenant_repo}
}

type DecommissionTenantRequest struct {
	TenantId string `json:"tenant_id"`
}

func (d DecommissionTenantCommand) Execute(req DecommissionTenantRequest) error {
	tenant_id, err := vo.NewTenantId(req.TenantId)
	if err != nil {
		return errors.BadRequest(8000, fmt.Sprintf("gagal melakukan dekomisi tenant: format id tenant tidak dikenal (ingin uuid, dapat %s)", req.TenantId))
	}
	tenant, err := d.tenant_repo.Find(tenant_id)
	if err != nil {
		return errors.Invariant(8001, fmt.Sprintf("terjadi kesalahan dalam mengambil data tenant: %s", err.Error()))
	}
	if tenant == nil {
		return errors.Invariant(8002, "data tenant tidak ditemukan")
	}
	tenant.Decommission()
	return d.tenant_repo.Persist(tenant)
}
