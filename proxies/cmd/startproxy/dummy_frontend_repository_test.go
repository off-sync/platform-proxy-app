package startproxy

import (
	"errors"

	"github.com/off-sync/platform-proxy-app/interfaces"
	"github.com/off-sync/platform-proxy-domain/frontends"
)

type dummyFrontendRepository struct {
	frontendNames []string
}

func (r *dummyFrontendRepository) FindAll() ([]*frontends.Frontend, error) {
	if len(r.frontendNames) < 1 {
		// return error in case the list is empty
		return nil, errors.New("no frontend URLs configured")
	}

	fs := make([]*frontends.Frontend, len(r.frontendNames))
	for i, n := range r.frontendNames {
		fs[i] = mockFrontend(n)
	}

	return fs, nil
}

func (r *dummyFrontendRepository) FindByName(name string) (*frontends.Frontend, error) {
	for _, n := range r.frontendNames {
		if name == n {
			return mockFrontend(n), nil
		}
	}

	return nil, interfaces.ErrUnknownFrontend
}

func (r *dummyFrontendRepository) Subscribe(events chan<- interfaces.FrontendEvent) {
	// do nothing
}

func mockFrontend(name string) *frontends.Frontend {
	f, err := frontends.NewFrontend(name, "http://"+name, nil, "service:"+name)
	if err != nil {
		// should not happen
		panic(err)
	}

	return f
}
