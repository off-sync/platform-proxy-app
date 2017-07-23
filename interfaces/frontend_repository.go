package interfaces

import (
	"errors"

	"github.com/off-sync/platform-proxy-domain/frontends"
)

// Errors
var (
	ErrUnknownFrontend = errors.New("unknown frontend")
)

// FrontendRepository is a repository for frontends.
type FrontendRepository interface {
	// FindAll returns all frontends contained in this repository.
	FindAll() ([]*frontends.Frontend, error)

	// FindByName returns the frontend with the specified name. If no frontend
	// exists with that name an ErrUnknownFrontend is returned.
	FindByName(name string) (*frontends.Frontend, error)
}
