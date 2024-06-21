package postgres

import (
	"encoding/json"
	"github.com/Marcellinom/tenant-management-saas/internal/app/queries"
	"github.com/Marcellinom/tenant-management-saas/provider"
)

type TenantQuery struct {
	db *provider.Database
}

func NewTenantQuery(db *provider.Database) *TenantQuery {
	return &TenantQuery{db: db}
}

func (t TenantQuery) GetByOrganizationId(orgs_id string) ([]queries.TenantQueryResult, error) {
	var res []struct {
		Id                  string
		Name                string
		Status              string
		ResourceInformation []byte
	}
	err := t.db.Table("tenants").
		Where("organization_id", orgs_id).
		Where("deleted_at is null").Find(&res).Error
	if err != nil {
		return nil, err
	}
	query_res := make([]queries.TenantQueryResult, len(res))
	for i, v := range res {
		var resource map[string]any
		if v.ResourceInformation != nil {
			json.Unmarshal(v.ResourceInformation, &resource)
		}
		query_res[i] = queries.TenantQueryResult{
			TenantId:            v.Id,
			Name:                v.Name,
			Status:              v.Status,
			ResourceInformation: resource,
		}
	}
	return query_res, nil
}
