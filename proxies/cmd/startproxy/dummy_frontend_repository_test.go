package startproxy

import (
	"errors"
	"strings"
	"time"

	"github.com/off-sync/platform-proxy-app/interfaces"
	"github.com/off-sync/platform-proxy-domain/frontends"
)

type dummyFrontendRepository struct {
	frontendNames []string
}

func (r *dummyFrontendRepository) ListFrontends() ([]string, error) {
	if len(r.frontendNames) < 1 {
		// return error in case the list is empty
		return nil, errors.New("no frontend URLs configured")
	}

	return r.frontendNames, nil
}

func (r *dummyFrontendRepository) DescribeFrontend(name string) (*frontends.Frontend, error) {
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
	var scheme = "http://"
	var cert *frontends.Certificate

	if strings.HasPrefix(name, "secure-") {
		scheme = "https://"
		cert = &frontends.Certificate{}
	}

	f, err := frontends.NewFrontend(name, scheme+name, cert, name)
	if err != nil {
		// should not happen
		panic(err)
	}

	return f
}
