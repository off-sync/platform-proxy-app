package getfrontends

import (
	"errors"
	"testing"

	"github.com/off-sync/platform-proxy-app/interfaces"
	"github.com/off-sync/platform-proxy-domain/frontends"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	q, err := New(&dummyRepo{})

	assert.NotNil(t, q)
	assert.Nil(t, err)
}

func TestNewShouldReturnErrorOnMissingRepository(t *testing.T) {
	q, err := New(nil)

	assert.Nil(t, q)
	assert.NotNil(t, err)

	assert.Equal(t, ErrMissingFrontendRepository, err)
}

func TestExecute(t *testing.T) {
	q, _ := New(&dummyRepo{
		frontendNames: []string{"test1", "test2"},
	})

	r, err := q.Execute(&Model{})

	assert.NotNil(t, r)
	assert.Nil(t, err)
}

func TestExecuteShouldReturnErrorFromRepository(t *testing.T) {
	q, _ := New(&dummyRepo{})

	r, err := q.Execute(&Model{})

	assert.Nil(t, r)
	assert.NotNil(t, err)
}

type dummyRepo struct {
	frontendNames []string
}

func (r *dummyRepo) FindAll() ([]*frontends.Frontend, error) {
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

func (r *dummyRepo) FindByName(name string) (*frontends.Frontend, error) {
	for _, n := range r.frontendNames {
		if name == n {
			return mockFrontend(n), nil
		}
	}

	return nil, interfaces.ErrUnknownFrontend
}

func mockFrontend(name string) *frontends.Frontend {
	f, err := frontends.NewFrontend(name, "http://"+name, nil, "service:"+name)
	if err != nil {
		// should not happen
		panic(err)
	}

	return f
}
