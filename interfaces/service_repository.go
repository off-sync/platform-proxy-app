package interfaces

import (
	"errors"

	"github.com/off-sync/platform-proxy-domain/services"
)

// Errors
var (
	ErrUnknownService = errors.New("unknown service")
)

// ServiceRepository is a repository for services.
type ServiceRepository interface {
	// FindAll returns all services contained in this repository.
	FindAll() ([]*services.Service, error)

	// FindByName returns the service with the specified name. If no service
	// exists with that name an ErrUnknownService is returned.
	FindByName(name string) (*services.Service, error)
}
