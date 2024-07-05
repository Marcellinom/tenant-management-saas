package queries

type TenantQueryInterface interface {
	GetByOrganizationId(organization_id string) ([]TenantQueryResult, error)
	FindByOrganizationAndAppId(organization_id string, app_id int) (*TenantQueryResult, error)
	Find(organization_id, tenant_id string) (*TenantQueryResult, error)
}

type TenantQueryResult struct {
	TenantId            string         `json:"tenant_id"`
	Name                string         `json:"name"`
	Status              string         `json:"status"`
	ResourceInformation map[string]any `json:"resource_information"`
	ProductId           string         `json:"product_id"`
	Tier                string         `json:"tier"`
	AppId               int            `json:"app_id"`
	AppName             string         `json:"app_name"`
}
