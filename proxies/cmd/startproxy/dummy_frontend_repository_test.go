package startproxy

import (
	"errors"
	"fmt"
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
	if name == "fail" {
		return nil, fmt.Errorf("DescribeFrontend(%s)", name)
	}

	for _, n := range r.frontendNames {
		if name == n {
			return mockFrontend(n), nil
		}
	}

	return nil, interfaces.ErrUnknownFrontend
}

func (r *dummyFrontendRepository) Subscribe() <-chan interfaces.FrontendEvent {
	events := make(chan interfaces.FrontendEvent)

	time.AfterFunc(250*time.Millisecond, func() {
		events <- interfaces.FrontendEvent{Name: "testapp"}
	})

	time.AfterFunc(350*time.Millisecond, func() {
		events <- interfaces.FrontendEvent{Name: "unknown"}
	})

	time.AfterFunc(450*time.Millisecond, func() {
		r.frontendNames = []string{}

		events <- interfaces.FrontendEvent{Name: "testapp"}
		events <- interfaces.FrontendEvent{Name: "secure-testapp"}
	})

	return events
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
