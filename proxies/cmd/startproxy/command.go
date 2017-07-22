package startproxy

import (
	"errors"

	"github.com/off-sync/platform-proxy-app/interfaces"
)

// Errors
var (
	ErrFrontendRepositoryMissing = errors.New("frontend repository missing")
	ErrServiceRepositoryMissing  = errors.New("service repository missing")
)

type Command struct {
	frontendRepository interfaces.FrontendRepository
	frontendWatcher    interfaces.FrontendWatcher
	serviceRepository  interfaces.ServiceRepository
	serviceWatcher     interfaces.ServiceWatcher
}

func NewCommand(
	frontendRepository interfaces.FrontendRepository,
	frontendWatcher interfaces.FrontendWatcher,
	serviceRepository interfaces.ServiceRepository,
	serviceWatcher interfaces.ServiceWatcher) (*Command, error) {
	if frontendRepository == nil {
		return nil, ErrFrontendRepositoryMissing
	}

	if serviceRepository == nil {
		return nil, ErrServiceRepositoryMissing
	}

	return &Command{
		frontendRepository: frontendRepository,
		frontendWatcher:    frontendWatcher,
		serviceRepository:  serviceRepository,
		serviceWatcher:     serviceWatcher,
	}, nil
}

func (c *Command) Execute(model *Model) error {
	return nil
}
