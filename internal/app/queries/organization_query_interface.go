package queries

type OrganizationQueryInterface interface {
	GetOrganizationByUserId(user_id string) ([]OrganizationResult, error)
}

type OrganizationResult struct {
	OrganizationId string `json:"organization_id"`
	Name           string `json:"name"`
}
