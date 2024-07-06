package listeners

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/events"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/repositories"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/vo"
	"github.com/Marcellinom/tenant-management-saas/provider/event"
)

type RegisteringTenantResource struct {
	tenant_repo repositories.TenantRepositoryInterface
}

func NewRegisteringTenantResource(tenant_repo repositories.TenantRepositoryInterface) *RegisteringTenantResource {
	return &RegisteringTenantResource{tenant_repo: tenant_repo}
}

func (l RegisteringTenantResource) Handle(ctx context.Context, event event.Event) error {
	var payload events.TenantResourceRegistered
	json_data, err := event.JSON()
	if err != nil {
		return fmt.Errorf("gagal menencode json pada event listener: %w", err)
	}
	err = json.Unmarshal(json_data, &payload)
	if err != nil {
		return fmt.Errorf("gagal mendecode json pada event listener: %w", err)
	}

	tenant_id, err := vo.NewTenantId(payload.TenantId)
	if err != nil {
		return err
	}
	tenant, err := l.tenant_repo.Find(tenant_id)
	if err != nil {
		return err
	}
	if tenant == nil {
		return fmt.Errorf("data tenant dengan id %s tidak ditemukan", payload.TenantId)
	}

	var metadata, resource_information []byte
	var resource string
	if r, ok := payload.ResourceInformation.(string); ok && r != "" {
		resource = r
	}
	if r, ok := payload.Metadata.(string); ok && r != "" {
		resource = r
	}

	metadata, err = base64.StdEncoding.DecodeString(resource)
	if err != nil {
		// kalo bukan b64 coba langsung encode jadi []byte
		metadata = []byte(resource)
	}
	if !json.Valid(metadata) {
		return fmt.Errorf("invalid json format saat registrasi tenant resource: %s", string(metadata))
	}

	var metadata_map map[string]any
	err = json.Unmarshal(metadata, &metadata_map)
	if err != nil {
		return fmt.Errorf("failed to decode metadata json: %w", err)
	}

	if v, exists := metadata_map["resource_information"]; exists {
		resource_information, _ = json.Marshal(v)
	}

	err = tenant.ActivateWithNewResourceInformation(resource_information)
	if err != nil {
		return fmt.Errorf("gagal melakukan registrasi resource tenant: %w", err)
	}
	return l.tenant_repo.Persist(tenant)
}

func (l RegisteringTenantResource) MaxRetries() int {
	return 3
}

func (l RegisteringTenantResource) Name() string {
	return fmt.Sprintf("%T", l)
}
