package startproxy

import (
	"context"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/off-sync/platform-proxy-app/infra/logging"
	"github.com/off-sync/platform-proxy-app/interfaces"
	"github.com/stretchr/testify/assert"
)

var logger interfaces.Logger

func init() {
	l := logrus.New()
	l.Level = logrus.DebugLevel

	logger = logging.NewLogrusLogger(l)
}

func TestNewCommand(t *testing.T) {
	c, err := NewCommand(
		&dummyFrontendRepository{},
		&dummyServiceRepository{},
		logger)

	assert.NotNil(t, c)
	assert.Nil(t, err)
}

func TestNewCommandShouldReturnErrorOnMissingFrontendRepository(t *testing.T) {
	c, err := NewCommand(
		nil,
		&dummyServiceRepository{},
		logger)

	assert.Nil(t, c)
	assert.NotNil(t, err)

	assert.Equal(t, ErrFrontendRepositoryMissing, err)
}

func TestNewCommandShouldReturnErrorOnMissingServiceRepository(t *testing.T) {
	c, err := NewCommand(
		&dummyFrontendRepository{},
		nil,
		logger)

	assert.Nil(t, c)
	assert.NotNil(t, err)

	assert.Equal(t, ErrServiceRepositoryMissing, err)
}

func TestNewCommandWithWatchers(t *testing.T) {
	fr := &dummyFrontendRepository{}
	sr := &dummyServiceRepository{}

	c, err := NewCommandWithWatchers(fr, sr, fr, sr, logger)

	assert.NotNil(t, c)
	assert.Nil(t, err)
}

func TestNewCommandWithWatchersShouldReturnErrorOnMissingFrontendRepository(t *testing.T) {
	fr := &dummyFrontendRepository{}
	sr := &dummyServiceRepository{}

	c, err := NewCommandWithWatchers(nil, sr, fr, sr, logger)

	assert.Nil(t, c)
	assert.NotNil(t, err)

	assert.Equal(t, ErrFrontendRepositoryMissing, err)
}

func TestNewCommandWithWatchersShouldReturnErrorOnMissingServiceRepository(t *testing.T) {
	fr := &dummyFrontendRepository{}
	sr := &dummyServiceRepository{}

	c, err := NewCommandWithWatchers(fr, nil, fr, sr, logger)

	assert.Nil(t, c)
	assert.NotNil(t, err)

	assert.Equal(t, ErrServiceRepositoryMissing, err)
}

func TestNewCommandWithWatchersShouldReturnErrorOnMissingFrontendWatcher(t *testing.T) {
	fr := &dummyFrontendRepository{}
	sr := &dummyServiceRepository{}

	c, err := NewCommandWithWatchers(fr, sr, nil, sr, logger)

	assert.Nil(t, c)
	assert.NotNil(t, err)

	assert.Equal(t, ErrFrontendWatcherMissing, err)
}

func TestNewCommandWithWatchersShouldReturnErrorOnMissingServiceWatcher(t *testing.T) {
	fr := &dummyFrontendRepository{}
	sr := &dummyServiceRepository{}

	c, err := NewCommandWithWatchers(fr, sr, fr, nil, logger)

	assert.Nil(t, c)
	assert.NotNil(t, err)

	assert.Equal(t, ErrServiceWatcherMissing, err)
}

func TestExecute(t *testing.T) {
	fr := &dummyFrontendRepository{frontendNames: []string{"testapp"}}
	sr := &dummyServiceRepository{serviceNames: []string{"testapp"}}

	c, _ := NewCommandWithWatchers(fr, sr, fr, sr, logger)

	ctx, cancel := context.WithCancel(context.Background())

	err := c.Execute(&Model{
		Ctx:             ctx,
		HTTPWebServer:   &dummyWebServer{},
		HTTPSWebServer:  &dummyWebServer{},
		LoadBalancer:    &dummyLoadBalancer{},
		PollingDuration: 1 * time.Second,
	})

	assert.Nil(t, err)

	time.Sleep(3 * time.Second)

	cancel()
}

func TestExecuteShouldReturnErrorOnMissingWebServers(t *testing.T) {
	fr := &dummyFrontendRepository{}
	sr := &dummyServiceRepository{}

	c, _ := NewCommandWithWatchers(fr, sr, fr, sr, logger)

	err := c.Execute(&Model{})

	assert.NotNil(t, err)

	assert.Equal(t, ErrWebServersMissing, err)
}

func TestExecuteShouldAcceptNilContext(t *testing.T) {
	fr := &dummyFrontendRepository{}
	sr := &dummyServiceRepository{}

	c, _ := NewCommandWithWatchers(fr, sr, fr, sr, logger)

	err := c.Execute(&Model{
		HTTPWebServer:   &dummyWebServer{},
		HTTPSWebServer:  &dummyWebServer{},
		Ctx:             nil,
		LoadBalancer:    &dummyLoadBalancer{},
		PollingDuration: 60 * time.Second,
	})

	assert.Nil(t, err)
}

func TestExecuteShouldReturnErrorOnNegativeDuration(t *testing.T) {
	fr := &dummyFrontendRepository{}
	sr := &dummyServiceRepository{}

	c, _ := NewCommandWithWatchers(fr, sr, fr, sr, logger)

	ctx, cancel := context.WithCancel(context.Background())

	err := c.Execute(&Model{
		HTTPWebServer:   &dummyWebServer{},
		HTTPSWebServer:  &dummyWebServer{},
		Ctx:             ctx,
		LoadBalancer:    &dummyLoadBalancer{},
		PollingDuration: -1 * time.Second,
	})

	assert.NotNil(t, err)

	assert.Equal(t, ErrInvalidPollingDuration, err)

	cancel()
}
