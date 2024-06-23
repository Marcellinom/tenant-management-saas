package events

import (
	"encoding/json"
	"time"
)

type TenantPersisted struct {
	TenantId         string    `json:"tenant_id"`
	ProductId        string    `json:"product_id"`
	OrganizationId   string    `json:"organization_id"`
	InfrastructureId string    `json:"infrastructure_id"`
	Name             string    `json:"name"`
	TenantStatus     string    `json:"tenant_status"`
	Timestamp        time.Time `json:"timestamp"`
}

func NewTenantPersisted(
	tenantId string,
	productId string,
	organizationId string,
	infrastructureId string,
	name string,
	tenantStatus string) *TenantPersisted {
	return &TenantPersisted{
		TenantId:         tenantId,
		ProductId:        productId,
		OrganizationId:   organizationId,
		InfrastructureId: infrastructureId,
		Name:             name,
		TenantStatus:     tenantStatus,
		Timestamp:        time.Now()}
}

func (t TenantPersisted) OccuredOn() time.Time {
	return t.Timestamp
}

func (t TenantPersisted) JSON() ([]byte, error) {
	return json.Marshal(t)
}
