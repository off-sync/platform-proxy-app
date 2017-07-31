package startproxy

import (
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/off-sync/platform-proxy-app/infra/logging"
	"github.com/off-sync/platform-proxy-app/interfaces"
)

var logger interfaces.Logger

func init() {
	l := logrus.New()
	l.Level = logrus.DebugLevel

	logger = logging.NewLogrusLogger(l)
}

func TestNewCommand(t *testing.T) {
	c, err := NewCommand(
		&dummyServiceRepository{},
		&dummyFrontendRepository{},
		logger)

	assert.NotNil(t, c)
	assert.Nil(t, err)
}

func TestNewCommandShouldReturnErrorOnMissingServiceRepository(t *testing.T) {
	c, err := NewCommand(
		nil,
		&dummyFrontendRepository{},
		logger)

	assert.Nil(t, c)
	assert.NotNil(t, err)

	assert.Equal(t, ErrServiceRepositoryMissing, err)
}

func TestNewCommandShouldReturnErrorOnMissingFrontendRepository(t *testing.T) {
	c, err := NewCommand(
		&dummyServiceRepository{},
		nil,
		logger)

	assert.Nil(t, c)
	assert.NotNil(t, err)

	assert.Equal(t, ErrFrontendRepositoryMissing, err)
}

func TestExecute(t *testing.T) {
	sr := &dummyServiceRepository{serviceNames: []string{"testapp"}}
	fr := &dummyFrontendRepository{frontendNames: []string{"testapp", "secure-testapp", "noservice"}}

	c, _ := NewCommand(sr, fr, logger)

	ctx, cancel := context.WithCancel(context.Background())

	err := c.Execute(&Model{
		Ctx:             ctx,
		WebServer:       &dummyWebServer{},
		SecureWebServer: &dummyWebServer{},
		LoadBalancer:    &dummyLoadBalancer{},
		PollingDuration: 500 * time.Millisecond,
	})

	assert.Nil(t, err)

	time.Sleep(1100 * time.Millisecond)

	cancel()
}

func TestExecuteShouldLogRepositoryErrors(t *testing.T) {
	sr := &dummyServiceRepository{}
	fr := &dummyFrontendRepository{}

	c, _ := NewCommand(sr, fr, logger)

	ctx, cancel := context.WithCancel(context.Background())

	err := c.Execute(&Model{
		Ctx:             ctx,
		WebServer:       &dummyWebServer{},
		SecureWebServer: &dummyWebServer{},
		LoadBalancer:    &dummyLoadBalancer{},
		PollingDuration: 60 * time.Second,
	})

	assert.Nil(t, err)

	cancel()
}

func TestExecuteShouldLogDescribeServiceErrors(t *testing.T) {
	sr := &dummyServiceRepository{[]string{"fail"}}
	fr := &dummyFrontendRepository{[]string{"testapp"}}

	c, _ := NewCommand(sr, fr, logger)

	ctx, cancel := context.WithCancel(context.Background())

	err := c.Execute(&Model{
		Ctx:             ctx,
		WebServer:       &dummyWebServer{},
		SecureWebServer: &dummyWebServer{},
		LoadBalancer:    &dummyLoadBalancer{},
		PollingDuration: 60 * time.Second,
	})

	assert.Nil(t, err)

	cancel()
}

func TestExecuteShouldLogDescribeFrontendErrors(t *testing.T) {
	sr := &dummyServiceRepository{[]string{"testapp"}}
	fr := &dummyFrontendRepository{[]string{"fail"}}

	c, _ := NewCommand(sr, fr, logger)

	ctx, cancel := context.WithCancel(context.Background())

	err := c.Execute(&Model{
		Ctx:             ctx,
		WebServer:       &dummyWebServer{},
		SecureWebServer: &dummyWebServer{},
		LoadBalancer:    &dummyLoadBalancer{},
		PollingDuration: 60 * time.Second,
	})

	assert.Nil(t, err)

	cancel()
}

func TestExecuteShouldLogWebServerErrors(t *testing.T) {
	sr := &dummyServiceRepository{serviceNames: []string{"testapp"}}
	fr := &dummyFrontendRepository{frontendNames: []string{"testapp", "secure-testapp"}}

	c, _ := NewCommand(sr, fr, logger)

	ctx, cancel := context.WithCancel(context.Background())

	err := c.Execute(&Model{
		Ctx:             ctx,
		WebServer:       &dummyWebServer{FailAll: true},
		SecureWebServer: &dummyWebServer{FailAll: true},
		LoadBalancer:    &dummyLoadBalancer{},
		PollingDuration: 60 * time.Second,
	})

	assert.Nil(t, err)

	cancel()
}

func TestExecuteShouldLogLoadBalancerErrors(t *testing.T) {
	sr := &dummyServiceRepository{serviceNames: []string{"testapp"}}
	fr := &dummyFrontendRepository{frontendNames: []string{"testapp", "testapp2"}}

	c, _ := NewCommand(sr, fr, logger)

	ctx, cancel := context.WithCancel(context.Background())

	web := &dummyWebServer{}

	err := c.Execute(&Model{
		Ctx:             ctx,
		WebServer:       web,
		SecureWebServer: web,
		LoadBalancer:    &dummyLoadBalancer{FailAll: true},
		PollingDuration: 60 * time.Second,
	})

	assert.Nil(t, err)

	// give proxy time to configure
	time.Sleep(200 * time.Millisecond)

	u, _ := url.Parse("http://testapp")
	resp := web.Handle(u, &http.Request{})

	assert.Equal(t, "Service not configured\n", resp)

	cancel()
}

func TestExecuteShouldReturnErrorOnMissingWebServer(t *testing.T) {
	sr := &dummyServiceRepository{}
	fr := &dummyFrontendRepository{}

	c, _ := NewCommand(sr, fr, logger)

	err := c.Execute(&Model{
		Ctx:             context.Background(),
		SecureWebServer: &dummyWebServer{},
		LoadBalancer:    &dummyLoadBalancer{},
	})

	assert.NotNil(t, err)

	assert.Equal(t, ErrWebServerMissing, err)
}

func TestExecuteShouldReturnErrorOnMissingSecureWebServer(t *testing.T) {
	sr := &dummyServiceRepository{}
	fr := &dummyFrontendRepository{}

	c, _ := NewCommand(sr, fr, logger)

	err := c.Execute(&Model{
		Ctx:          context.Background(),
		WebServer:    &dummyWebServer{},
		LoadBalancer: &dummyLoadBalancer{},
	})

	assert.NotNil(t, err)

	assert.Equal(t, ErrSecureWebServerMissing, err)
}

func TestExecuteShouldReturnErrorOnMissingLoadBalancer(t *testing.T) {
	sr := &dummyServiceRepository{}
	fr := &dummyFrontendRepository{}

	c, _ := NewCommand(sr, fr, logger)

	err := c.Execute(&Model{
		Ctx:             context.Background(),
		WebServer:       &dummyWebServer{},
		SecureWebServer: &dummyWebServer{},
	})

	assert.NotNil(t, err)

	assert.Equal(t, ErrLoadBalancerMissing, err)
}

func TestExecuteShouldAcceptNilContext(t *testing.T) {
	sr := &dummyServiceRepository{}
	fr := &dummyFrontendRepository{}

	c, _ := NewCommand(sr, fr, logger)

	err := c.Execute(&Model{
		Ctx:             nil,
		WebServer:       &dummyWebServer{},
		SecureWebServer: &dummyWebServer{},
		LoadBalancer:    &dummyLoadBalancer{},
		PollingDuration: 60 * time.Second,
	})

	assert.Nil(t, err)
}

func TestExecuteShouldReturnErrorOnNegativeDuration(t *testing.T) {
	sr := &dummyServiceRepository{}
	fr := &dummyFrontendRepository{}

	c, _ := NewCommand(sr, fr, logger)

	ctx, cancel := context.WithCancel(context.Background())

	err := c.Execute(&Model{
		Ctx:             ctx,
		WebServer:       &dummyWebServer{},
		SecureWebServer: &dummyWebServer{},
		LoadBalancer:    &dummyLoadBalancer{},
		PollingDuration: -1 * time.Second,
	})

	assert.NotNil(t, err)

	assert.Equal(t, ErrInvalidPollingDuration, err)

	cancel()
}
