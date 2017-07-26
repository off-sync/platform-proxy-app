package startproxy

import (
	"errors"
	"time"

	"github.com/off-sync/platform-proxy-app/interfaces"
	"github.com/off-sync/platform-proxy-domain/services"
)

type dummyServiceRepository struct {
	serviceNames []string
}

func (r *dummyServiceRepository) FindAll() ([]*services.Service, error) {
	if len(r.serviceNames) < 1 {
		// return error in case the list is empty
		return nil, errors.New("no frontend URLs configured")
	}

	fs := make([]*services.Service, len(r.serviceNames))
	for i, n := range r.serviceNames {
		fs[i] = mockService(n)
	}

	return fs, nil
}

func (r *dummyServiceRepository) FindByName(name string) (*services.Service, error) {
	for _, n := range r.serviceNames {
		if name == n {
			return mockService(n), nil
		}
	}

	return nil, interfaces.ErrUnknownFrontend
}

func (r *dummyServiceRepository) Subscribe(events chan<- interfaces.ServiceEvent) {
	time.AfterFunc(2200*time.Millisecond, func() {
		events <- interfaces.ServiceEvent{Name: "testapp"}
	})
}

func mockService(name string) *services.Service {
	f, err := services.NewService(name, "http://127.0.0.1:8080", "http://127.0.0.1:8080")
	if err != nil {
		// should not happen
		panic(err)
	}

	return f
}
