package startproxy

import (
	"context"
	"errors"

	"github.com/off-sync/platform-proxy-app/interfaces"
)

// Errors
var (
	ErrFrontendRepositoryMissing = errors.New("frontend repository missing")
	ErrFrontendWatcherMissing    = errors.New("frontend watcher missing")
	ErrServiceRepositoryMissing  = errors.New("service repository missing")
	ErrServiceWatcherMissing     = errors.New("service watcher missing")
	ErrWebServersMissing         = errors.New("both web servers missing, provide at least one")
	ErrInvalidPollingDuration    = errors.New("invalid polling duration, must greater than or equal to 0")
)

// Command models the Start Proxy Command which can be used to start one of the
// platform proxies.
type Command struct {
	frontendRepository interfaces.FrontendRepository
	frontendWatcher    interfaces.FrontendWatcher
	serviceRepository  interfaces.ServiceRepository
	serviceWatcher     interfaces.ServiceWatcher
}

// NewCommand creates a new Start Proxy Command using the provided frontend
// and service repositories.
func NewCommand(
	frontendRepository interfaces.FrontendRepository,
	serviceRepository interfaces.ServiceRepository) (*Command, error) {
	if frontendRepository == nil {
		return nil, ErrFrontendRepositoryMissing
	}

	if serviceRepository == nil {
		return nil, ErrServiceRepositoryMissing
	}

	return &Command{
		frontendRepository: frontendRepository,
		serviceRepository:  serviceRepository,
	}, nil
}

// NewCommandWithWatchers creates a new Start Proxy Command including watchers
// that will be used to update the frontends and services.
func NewCommandWithWatchers(frontendRepository interfaces.FrontendRepository,
	serviceRepository interfaces.ServiceRepository,
	frontendWatcher interfaces.FrontendWatcher,
	serviceWatcher interfaces.ServiceWatcher) (*Command, error) {
	c, err := NewCommand(frontendRepository, serviceRepository)
	if err != nil {
		return nil, err
	}

	if frontendWatcher == nil {
		return nil, ErrFrontendWatcherMissing
	}

	if serviceWatcher == nil {
		return nil, ErrServiceWatcherMissing
	}

	c.frontendWatcher = frontendWatcher
	c.serviceWatcher = serviceWatcher

	return c, nil
}

// Execute runs the Start Proxy Command by configuring the required listeners.
func (c *Command) Execute(model *Model) error {
	if model.HTTPWebServer == nil && model.HTTPSWebServer == nil {
		return ErrWebServersMissing
	}

	if model.PollingDuration < 0 {
		return ErrInvalidPollingDuration
	}

	if model.Ctx == nil {
		model.Ctx = context.Background()
	}

	go runProxy(c, model)

	return nil
}
