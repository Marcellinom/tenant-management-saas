package iam

import (
	"github.com/Marcellinom/tenant-management-saas/internal/app/queries"
	"github.com/Marcellinom/tenant-management-saas/provider"
)

type OrganizationQuery struct {
	db *provider.Database
}

func NewOrganizationQuery(db *provider.Database) *OrganizationQuery {
	return &OrganizationQuery{db: db}
}

func (q OrganizationQuery) GetOrganizationByUserId(user_id string) ([]queries.OrganizationResult, error) {
	var Orgs []struct {
		Id   string
		Name string
	}
	err := q.db.Table("organization o").
		Joins("join user_organization uo on uo.user_id = ? and uo.organization_id = o.id", user_id).
		Find(&Orgs).Error
	if err != nil {
		return nil, err
	}
	res := make([]queries.OrganizationResult, len(Orgs))
	for i, v := range Orgs {
		res[i] = queries.OrganizationResult{
			OrganizationId: v.Id,
			Name:           v.Name,
		}
	}
	return res, nil
}
