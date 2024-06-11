package services

import (
	"github.com/Marcellinom/tenant-management-saas/internal/domain/entities/Infrastructure"
	"github.com/google/uuid"
)

type InfrastructureServiceInterface interface {
	FindAvailablePool() (*Infrastructure.Infrastructure, error)
	Persist(infra *Infrastructure.Infrastructure) error
	Find(infra_id uuid.UUID) (*Infrastructure.Infrastructure, error)
}
