package postgres

import (
	"encoding/json"
	"errors"
	"github.com/Marcellinom/tenant-management-saas/internal/app/queries"
	"github.com/Marcellinom/tenant-management-saas/provider"
	"gorm.io/gorm"
)

type TenantQuery struct {
	db *provider.Database
}

func NewTenantQuery(db *provider.Database) *TenantQuery {
	return &TenantQuery{db: db}
}

type tenantQueryDto struct {
	Id                  string
	Name                string
	Status              string
	ResourceInformation []byte
	ProductId           string
	Tier                string
	AppId               int
	AppName             string
}

func (t TenantQuery) GetByOrganizationId(orgs_id string) ([]queries.TenantQueryResult, error) {
	var res []tenantQueryDto
	err := t.db.Raw(
		"select "+
			"t.id,"+
			"name,"+
			"status,"+
			"resource_information,"+
			"p.id as product_id,"+
			"p.tier_name as tier,"+
			"p.app_id,"+
			"(select name from apps a where a.id = p.app_id) app_name "+
			"from (select id, name, status, resource_information, product_id from tenants where organization_id = ? and deleted_at is null"+
			") t join products p on t.product_id = p.id", orgs_id).Find(&res).Error
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
			ProductId:           v.ProductId,
			Tier:                v.Tier,
			AppId:               v.AppId,
			AppName:             v.AppName,
		}
	}
	return query_res, nil
}

func (t TenantQuery) Find(organization_id, tenant_id string) (*queries.TenantQueryResult, error) {
	var res tenantQueryDto
	err := t.db.Raw(
		"select "+
			"t.id,"+
			"name,"+
			"status,"+
			"resource_information,"+
			"p.id as product_id,"+
			"p.tier_name as tier,"+
			"p.app_id,"+
			"(select name from apps a where a.id = p.app_id) app_name "+
			"from (select id, name, status, resource_information, product_id from tenants where organization_id = ? and id = ? and deleted_at is null"+
			") t join products p on t.product_id = p.id", organization_id, tenant_id).
		Take(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}
	var resource map[string]any
	if res.ResourceInformation != nil {
		json.Unmarshal(res.ResourceInformation, &resource)
	}
	return &queries.TenantQueryResult{
		TenantId:            res.Id,
		Name:                res.Name,
		Status:              res.Status,
		ResourceInformation: resource,
		ProductId:           res.ProductId,
		Tier:                res.Tier,
		AppId:               res.AppId,
		AppName:             res.AppName,
	}, nil
}
