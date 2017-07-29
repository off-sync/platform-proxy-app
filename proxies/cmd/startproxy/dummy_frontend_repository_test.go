package startproxy

import (
	"errors"
	"time"

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
	time.AfterFunc(1200*time.Millisecond, func() {
		events <- interfaces.FrontendEvent{Name: "testapp"}
	})

	time.AfterFunc(1400*time.Millisecond, func() {
		events <- interfaces.FrontendEvent{Name: "unknown"}
	})
}

func mockFrontend(name string) *frontends.Frontend {
	f, err := frontends.NewFrontend(name, "http://"+name, nil, name)
	if err != nil {
		// should not happen
		panic(err)
	}

	return f
}
