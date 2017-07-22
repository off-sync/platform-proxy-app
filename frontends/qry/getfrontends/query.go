package getfrontends

import (
	"errors"

	"github.com/off-sync/platform-proxy-app/interfaces"
)

// Errors
var (
	ErrMissingFrontendRepository = errors.New("missing frontends repository")
)

// Query implements the Get Frontends Query. It requires a FrontendRepository.
type Query struct {
	repo interfaces.FrontendRepository
}

// New creates a new Get Frontends Query
func New(repo interfaces.FrontendRepository) (*Query, error) {
	if repo == nil {
		return nil, ErrMissingFrontendRepository
	}

	return &Query{
		repo: repo,
	}, nil
}

// Execute performs the Get Frontends Query using the provided model.
func (q *Query) Execute(model *Model) (*Result, error) {
	fs, err := q.repo.FindAll()
	if err != nil {
		return nil, err
	}

	return &Result{
		Frontends: fs,
	}, nil
}
