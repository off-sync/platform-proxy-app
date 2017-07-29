package getfrontends

import (
	"errors"
	"testing"

	"github.com/off-sync/platform-proxy-app/interfaces"
	"github.com/off-sync/platform-proxy-domain/frontends"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	q, err := NewQuery(&dummyRepo{})

	assert.NotNil(t, q)
	assert.Nil(t, err)
}

func TestNewShouldReturnErrorOnMissingRepository(t *testing.T) {
	q, err := NewQuery(nil)

	assert.Nil(t, q)
	assert.NotNil(t, err)

	assert.Equal(t, ErrMissingFrontendRepository, err)
}

func TestExecute(t *testing.T) {
	q, _ := NewQuery(&dummyRepo{
		frontendNames: []string{"test1", "test2"},
	})

	r, err := q.Execute(&Model{})

	assert.NotNil(t, r)
	assert.Nil(t, err)
}

func TestExecuteDescribeFrontendErrorShouldBeReturned(t *testing.T) {
	q, _ := NewQuery(&dummyRepo{
		frontendNames: []string{"unknown"},
	})

	r, err := q.Execute(&Model{})

	assert.Nil(t, r)
	assert.Equal(t, interfaces.ErrUnknownFrontend, err)
}

func TestExecuteShouldReturnErrorFromRepository(t *testing.T) {
	q, _ := NewQuery(&dummyRepo{})

	r, err := q.Execute(&Model{})

	assert.Nil(t, r)
	assert.NotNil(t, err)
}

type dummyRepo struct {
	frontendNames []string
}

func (r *dummyRepo) ListFrontends() ([]string, error) {
	if len(r.frontendNames) < 1 {
		// return error in case the list is empty
		return nil, errors.New("no frontend URLs configured")
	}

	return r.frontendNames, nil
}

func (r *dummyRepo) DescribeFrontend(name string) (*frontends.Frontend, error) {
	if name == "unknown" {
		return nil, interfaces.ErrUnknownFrontend
	}

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
