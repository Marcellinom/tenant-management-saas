package queries

type TenantQueryInterface interface {
	GetByOrganizationId(organization_id string) ([]TenantQueryResult, error)
}

type TenantQueryResult struct {
	TenantId            string         `json:"tenant_id"`
	Name                string         `json:"name"`
	Status              string         `json:"status"`
	ResourceInformation map[string]any `json:"resource_information"`
}
