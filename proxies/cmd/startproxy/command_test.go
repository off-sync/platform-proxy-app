package startproxy

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCommand(t *testing.T) {
	c, err := NewCommand(
		&dummyFrontendRepository{},
		&dummyServiceRepository{})

	assert.NotNil(t, c)
	assert.Nil(t, err)
}

func TestNewCommandShouldReturnErrorOnMissingFrontendRepository(t *testing.T) {
	c, err := NewCommand(
		nil,
		&dummyServiceRepository{})

	assert.Nil(t, c)
	assert.NotNil(t, err)

	assert.Equal(t, ErrFrontendRepositoryMissing, err)
}

func TestNewCommandShouldReturnErrorOnMissingServiceRepository(t *testing.T) {
	c, err := NewCommand(
		&dummyFrontendRepository{},
		nil)

	assert.Nil(t, c)
	assert.NotNil(t, err)

	assert.Equal(t, ErrServiceRepositoryMissing, err)
}

func TestNewCommandWithWatchers(t *testing.T) {
	fr := &dummyFrontendRepository{}
	sr := &dummyServiceRepository{}

	c, err := NewCommandWithWatchers(fr, sr, fr, sr)

	assert.NotNil(t, c)
	assert.Nil(t, err)
}

func TestNewCommandWithWatchersShouldReturnErrorOnMissingFrontendRepository(t *testing.T) {
	fr := &dummyFrontendRepository{}
	sr := &dummyServiceRepository{}

	c, err := NewCommandWithWatchers(nil, sr, fr, sr)

	assert.Nil(t, c)
	assert.NotNil(t, err)

	assert.Equal(t, ErrFrontendRepositoryMissing, err)
}

func TestNewCommandWithWatchersShouldReturnErrorOnMissingServiceRepository(t *testing.T) {
	fr := &dummyFrontendRepository{}
	sr := &dummyServiceRepository{}

	c, err := NewCommandWithWatchers(fr, nil, fr, sr)

	assert.Nil(t, c)
	assert.NotNil(t, err)

	assert.Equal(t, ErrServiceRepositoryMissing, err)
}

func TestNewCommandWithWatchersShouldReturnErrorOnMissingFrontendWatcher(t *testing.T) {
	fr := &dummyFrontendRepository{}
	sr := &dummyServiceRepository{}

	c, err := NewCommandWithWatchers(fr, sr, nil, sr)

	assert.Nil(t, c)
	assert.NotNil(t, err)

	assert.Equal(t, ErrFrontendWatcherMissing, err)
}

func TestNewCommandWithWatchersShouldReturnErrorOnMissingServiceWatcher(t *testing.T) {
	fr := &dummyFrontendRepository{}
	sr := &dummyServiceRepository{}

	c, err := NewCommandWithWatchers(fr, sr, fr, nil)

	assert.Nil(t, c)
	assert.NotNil(t, err)

	assert.Equal(t, ErrServiceWatcherMissing, err)
}

func TestExecute(t *testing.T) {
	fr := &dummyFrontendRepository{}
	sr := &dummyServiceRepository{}

	c, _ := NewCommandWithWatchers(fr, sr, fr, sr)

	ctx, cancel := context.WithCancel(context.Background())

	err := c.Execute(&Model{
		HTTPWebServer:   &dummyWebServer{},
		HTTPSWebServer:  &dummyWebServer{},
		Ctx:             ctx,
		PollingDuration: 60 * time.Second,
	})

	assert.Nil(t, err)

	cancel()
}

func TestExecuteShouldReturnErrorOnMissingWebServers(t *testing.T) {
	fr := &dummyFrontendRepository{}
	sr := &dummyServiceRepository{}

	c, _ := NewCommandWithWatchers(fr, sr, fr, sr)

	err := c.Execute(&Model{})

	assert.NotNil(t, err)

	assert.Equal(t, ErrWebServersMissing, err)
}

func TestExecuteShouldAcceptNilContext(t *testing.T) {
	fr := &dummyFrontendRepository{}
	sr := &dummyServiceRepository{}

	c, _ := NewCommandWithWatchers(fr, sr, fr, sr)

	err := c.Execute(&Model{
		HTTPWebServer:   &dummyWebServer{},
		HTTPSWebServer:  &dummyWebServer{},
		Ctx:             nil,
		PollingDuration: 60 * time.Second,
	})

	assert.Nil(t, err)
}

func TestExecuteShouldReturnErrorOnNegativeDuration(t *testing.T) {
	fr := &dummyFrontendRepository{}
	sr := &dummyServiceRepository{}

	c, _ := NewCommandWithWatchers(fr, sr, fr, sr)

	ctx, cancel := context.WithCancel(context.Background())

	err := c.Execute(&Model{
		HTTPWebServer:   &dummyWebServer{},
		HTTPSWebServer:  &dummyWebServer{},
		Ctx:             ctx,
		PollingDuration: -1 * time.Second,
	})

	assert.NotNil(t, err)

	assert.Equal(t, ErrInvalidPollingDuration, err)

	cancel()
}
