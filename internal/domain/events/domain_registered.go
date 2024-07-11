package events

import (
	"encoding/json"
	"time"
)

type DomainRegistered struct {
	AppId          int       `json:"app_id"`
	TenantId       string    `json:"tenant_id"`
	OrganizationId string    `json:"org_id"`
	DomainUrl      string    `json:"url"`
	Timestamp      time.Time `json:"timestamp"`
}

func NewDomainRegistered(appId int, tenantId string, organizationId string, domainUrl string) DomainRegistered {
	return DomainRegistered{AppId: appId, TenantId: tenantId, OrganizationId: organizationId, DomainUrl: domainUrl, Timestamp: time.Now()}
}

func (t DomainRegistered) OccuredOn() time.Time {
	return t.Timestamp
}

func (t DomainRegistered) JSON() ([]byte, error) {
	return json.Marshal(t)
}
