package startproxy

import (
	"context"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/off-sync/platform-proxy-app/interfaces"
)

type proxy struct {
	sync.Mutex

	// context
	ctx context.Context

	// logging
	logger interfaces.Logger

	// configuration
	frontendRepository interfaces.FrontendRepository
	serviceRepository  interfaces.ServiceRepository
	pollingDuration    time.Duration

	// request handling
	webServer       interfaces.WebServer
	secureWebServer interfaces.SecureWebServer
	loadBalancer    interfaces.LoadBalancer

	// internal state
	serviceHandlers map[string]http.Handler
}

func newProxy(
	ctx context.Context,
	logger interfaces.Logger,
	serviceRepository interfaces.ServiceRepository,
	frontendRepository interfaces.FrontendRepository,
	pollingDuration time.Duration,
	webServer interfaces.WebServer,
	secureWebServer interfaces.SecureWebServer,
	loadBalancer interfaces.LoadBalancer) *proxy {

	return &proxy{
		ctx:                ctx,
		logger:             logger,
		serviceRepository:  serviceRepository,
		frontendRepository: frontendRepository,
		pollingDuration:    pollingDuration,
		webServer:          webServer,
		secureWebServer:    secureWebServer,
		loadBalancer:       loadBalancer,
		serviceHandlers:    make(map[string]http.Handler),
	}
}

func (p *proxy) run() {
	// configure all services and frontends
	p.configure()

	// subscribe to service events
	serviceEvents := make(chan interfaces.ServiceEvent, 10)
	defer close(serviceEvents)

	if w, ok := p.serviceRepository.(interfaces.ServiceWatcher); ok {
		p.logger.Info("subscribing to service watcher")

		w.Subscribe(serviceEvents)
	}

	// subscribe to frontend events
	frontendEvents := make(chan interfaces.FrontendEvent, 10)
	defer close(frontendEvents)

	if w, ok := p.frontendRepository.(interfaces.FrontendWatcher); ok {
		p.logger.Info("subscribing to frontend watcher")

		w.Subscribe(frontendEvents)
	}

	// create polling ticker
	poll := time.Tick(p.pollingDuration)

	for {
		select {
		// respond to the context closing
		case <-p.ctx.Done():
			p.logger.Info("context is done: returning")
			return

		// respond to polling events
		case <-poll:
			p.logger.Info("polling configuration")
			p.configure()
			break

		// respond to service events
		case serviceEvent := <-serviceEvents:
			p.logger.
				WithField("name", serviceEvent.Name).
				Info("received service event")

			p.configureService(serviceEvent.Name)

			break

		// respond to frontend events
		case frontendEvent := <-frontendEvents:
			p.logger.
				WithField("name", frontendEvent.Name).
				Info("received frontend event")

			p.configureFrontend(frontendEvent.Name)

			break
		}
	}
}

func (p *proxy) configure() {
	// configure services first to create the required handlers
	services, err := p.serviceRepository.ListServices()
	if err != nil {
		p.logger.
			WithError(err).
			Error("listing services")
	} else {
		for _, service := range services {
			p.configureService(service)
		}
	}

	// configure frontends
	frontends, err := p.frontendRepository.ListFrontends()
	if err != nil {
		p.logger.
			WithError(err).
			Error("listing frontends")
	} else {
		for _, frontend := range frontends {
			p.configureFrontend(frontend)
		}
	}
}

func (p *proxy) getServiceHandler(serviceName string) http.Handler {
	handler, exists := p.serviceHandlers[serviceName]
	if !exists {
		return http.NotFoundHandler()
	}

	return handler
}

func (p *proxy) configureService(name string) {
	// describe service
	service, err := p.serviceRepository.DescribeService(name)
	if err != nil {
		p.logger.
			WithError(err).
			WithField("name", name).
			Error("describing service")

		return
	}

	// lock internal state
	p.Lock()
	defer p.Unlock()

	p.logger.
		WithField("name", service.Name).
		WithField("servers", service.Servers).
		Debug("configuring service")

	handler, err := p.loadBalancer.UpsertService(service.Name, service.Servers...)
	if err != nil {
		p.logger.
			WithError(err).
			WithField("name", name).
			Error("upserting service")

		// set the service handler to return an internal server error on each
		// request
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "service not configured", http.StatusInternalServerError)
		})
	}

	p.serviceHandlers[service.Name] = handler
}

func (p *proxy) configureFrontend(name string) {
	// get frontend from repository
	frontend, err := p.frontendRepository.DescribeFrontend(name)
	if err != nil {
		p.logger.
			WithError(err).
			WithField("name", name).
			Error("describing frontend")

		return
	}

	p.Lock()
	defer p.Unlock()

	p.logger.
		WithField("name", frontend.Name).
		WithField("url", frontend.URL).
		WithField("service_name", frontend.ServiceName).
		Debug("configuring frontend")

	if frontend.Certificate != nil {
		// configure HTTPS
		err := p.secureWebServer.UpsertCertificate(
			frontend.URL.Host,
			frontend.Certificate)
		if err != nil {
			p.logger.
				WithError(err).
				WithField("host", frontend.URL.Host).
				Error("upserting certificate")
		}

		err = p.secureWebServer.UpsertRoute(
			frontend.URL,
			p.getServiceHandler(frontend.ServiceName))
		if err != nil {
			p.logger.
				WithError(err).
				WithField("url", frontend.URL).
				Error("upserting route")
		}

		// configure HTTP redirect
		httpURL := &url.URL{}
		*httpURL = *frontend.URL
		httpURL.Scheme = "http"

		err = p.webServer.UpsertRoute(httpURL,
			http.RedirectHandler(
				frontend.URL.String(),
				http.StatusMovedPermanently))
		if err != nil {
			p.logger.
				WithError(err).
				WithField("url", frontend.URL).
				Error("upserting route")
		}
	} else {
		// configure HTTP
		err := p.webServer.UpsertRoute(
			frontend.URL,
			p.getServiceHandler(frontend.ServiceName))
		if err != nil {
			p.logger.
				WithError(err).
				WithField("url", frontend.URL).
				Error("upserting route")
		}
	}
}
