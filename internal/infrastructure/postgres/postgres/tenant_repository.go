package postgres

import (
	"errors"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/Tenant"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/events"
	"github.com/Marcellinom/tenant-management-saas/internal/domain/vo"
	"github.com/Marcellinom/tenant-management-saas/provider"
	"github.com/Marcellinom/tenant-management-saas/provider/event"
	"gorm.io/gorm"
	"time"
)

type TenantRepository struct {
	db            *provider.Database
	event_service event.Service
}

func NewTenantRepository(db *provider.Database, event_service event.Service) *TenantRepository {
	return &TenantRepository{db: db, event_service: event_service}
}

func (t TenantRepository) Find(tenant_id vo.TenantId) (*Tenant.Tenant, error) {
	var tenant_data struct {
		Id                  string
		ProductId           string
		OrganizationId      string
		InfrastructureId    string
		Name                string
		Status              string
		ResourceInformation []byte
	}
	err := t.db.Table("tenants").Where("id", tenant_id.String()).
		Take(&tenant_data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	id, _ := vo.NewTenantId(tenant_data.Id)
	productId, _ := vo.NewProductId(tenant_data.ProductId)
	organizationId, _ := vo.NewOrganizationId(tenant_data.OrganizationId)
	infrastructureId, _ := vo.NewInfrastructureId(tenant_data.InfrastructureId)

	return &Tenant.Tenant{
		TenantId:            id,
		ProductId:           productId,
		OrganizationId:      organizationId,
		TenantStatus:        Tenant.NewTenantStatus(Tenant.Status(tenant_data.Status)),
		Name:                tenant_data.Name,
		InfrastructureId:    infrastructureId,
		ResourceInformation: tenant_data.ResourceInformation,
		Events:              make(map[string]event.Event),
	}, nil
}

func (t TenantRepository) Persist(tenant *Tenant.Tenant) error {
	return t.db.Transaction(func(tx *gorm.DB) error {

		// event yang dihandle sebelum main persistance
		for _, e := range tenant.Events {
			switch e.(type) {
			case events.TenantDecommissioned:
				tx.Table("tenants").Where("id", tenant.TenantId.String()).
					Updates(map[string]any{
						"infrastructure_id": nil,
						"deleted_at":        time.Now(),
					})
				t.event_service.Dispatch(events.INFRASTRUCTURE_DESTROYED, events.NewInfrastructureDestroyed(
					tenant.InfrastructureId.String(),
				))
				return nil
			}
		}

		var row int64
		err := tx.Table("tenants").Where("id", tenant.TenantId.String()).
			Count(&row).Error
		if err != nil {
			return err
		}
		payload := map[string]any{
			"product_id":           tenant.ProductId.String(),
			"name":                 tenant.Name,
			"status":               tenant.TenantStatus,
			"updated_at":           time.Now(),
			"resource_information": tenant.ResourceInformation,
		}

		if row > 0 {
			payload["infrastructure_id"] = tenant.InfrastructureId.String()
			err = tx.Table("tenants").
				Where("id", tenant.TenantId.String()).
				Updates(payload).Error
		} else {
			payload["id"] = tenant.TenantId.String()
			payload["organization_id"] = tenant.OrganizationId.String()
			payload["created_at"] = time.Now()
			err = tx.Table("tenants").Create(payload).Error
		}
		if err == nil {
			t.event_service.Dispatch(events.TENANT_PERSISTED, events.NewTenantPersisted(
				tenant.TenantId.String(),
				tenant.ProductId.String(),
				tenant.OrganizationId.String(),
				tenant.InfrastructureId.String(),
				tenant.Name,
				tenant.TenantStatus,
			))
			// mesti dilakukan terakhir
			// supaya jika ada gagal diatas, gk kepublish
			for event_name, e := range tenant.Events {
				t.event_service.Dispatch(event_name, e)
			}
		}
		return err
	})
}
