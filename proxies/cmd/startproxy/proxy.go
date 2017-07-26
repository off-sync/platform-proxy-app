package startproxy

import (
	"context"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/off-sync/platform-proxy-app/interfaces"
	"github.com/off-sync/platform-proxy-domain/frontends"
	"github.com/off-sync/platform-proxy-domain/services"
)

type proxy struct {
	sync.Mutex

	// context
	ctx context.Context

	// logging
	logger interfaces.Logger

	// configuration
	frontendRepository interfaces.FrontendRepository
	frontendWatcher    interfaces.FrontendWatcher
	serviceRepository  interfaces.ServiceRepository
	serviceWatcher     interfaces.ServiceWatcher
	pollingDuration    time.Duration

	// web servers
	httpServer  interfaces.WebServer
	httpsServer interfaces.WebServer

	// internal state
	serviceHandlers map[string]http.Handler
}

func newProxy(
	ctx context.Context,
	logger interfaces.Logger,
	frontendRepository interfaces.FrontendRepository,
	frontendWatcher interfaces.FrontendWatcher,
	serviceRepository interfaces.ServiceRepository,
	serviceWatcher interfaces.ServiceWatcher,
	pollingDuration time.Duration,
	httpServer, httpsServer interfaces.WebServer) *proxy {
	return &proxy{
		ctx:                ctx,
		logger:             logger,
		frontendRepository: frontendRepository,
		frontendWatcher:    frontendWatcher,
		serviceRepository:  serviceRepository,
		serviceWatcher:     serviceWatcher,
		pollingDuration:    pollingDuration,
		httpServer:         httpServer,
		httpsServer:        httpsServer,
		serviceHandlers:    make(map[string]http.Handler),
	}
}

func (p *proxy) run() {
	services, err := p.serviceRepository.FindAll()
	if err != nil {
		p.logger.
			WithError(err).
			Error("unable to get services")
	} else {
		for _, service := range services {
			err = p.configureService(service)
			if err != nil {
				p.logger.
					WithError(err).
					Error("unable to configure service")
			}
		}
	}

	frontends, err := p.frontendRepository.FindAll()
	if err != nil {
		p.logger.
			WithError(err).
			Error("unable to get frontends")
	} else {
		for _, frontend := range frontends {
			err = p.configureFrontend(frontend)
			if err != nil {
				p.logger.
					WithError(err).
					Error("unable to configure frontend")
			}
		}
	}

	frontendEvents := make(chan interfaces.FrontendEvent, 10)
	defer close(frontendEvents)

	if p.frontendWatcher != nil {
		p.logger.Info("subscribing to frontend watcher")

		p.frontendWatcher.Subscribe(frontendEvents)
	}

	serviceEvents := make(chan interfaces.ServiceEvent, 10)
	defer close(serviceEvents)

	if p.serviceWatcher != nil {
		p.logger.Info("subscribing to service watcher")

		p.serviceWatcher.Subscribe(serviceEvents)
	}

	// create polling ticker
	if p.pollingDuration < 1 {
		p.pollingDuration = 5 * time.Minute
	}

	poll := time.Tick(p.pollingDuration)

	for {
		select {
		case <-p.ctx.Done():
			p.logger.Info("context is done: returning")
			return
		case <-poll:
			p.logger.Info("polling configuration")
			break
		case serviceEvent := <-serviceEvents:
			p.logger.
				WithField("name", serviceEvent.Name).
				Info("received service event")

			service, err := p.serviceRepository.FindByName(serviceEvent.Name)
			if err != nil {
				p.logger.
					WithError(err).
					Error("unable to find service by name")
			} else {
				p.configureService(service)
			}

			break
		case frontendEvent := <-frontendEvents:
			p.logger.
				WithField("name", frontendEvent.Name).
				Info("received frontend event")

			frontend, err := p.frontendRepository.FindByName(frontendEvent.Name)
			if err != nil {
				p.logger.
					WithError(err).
					Error("unable to find frontend by name")
			} else {
				p.configureFrontend(frontend)
			}

			break
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

func (p *proxy) configureService(service *services.Service) error {
	p.Lock()
	defer p.Unlock()

	p.logger.
		WithField("name", service.Name).
		WithField("servers", service.Servers).
		Debug("configuring service")

	return nil
}

func (p *proxy) configureFrontend(frontend *frontends.Frontend) error {
	p.Lock()
	defer p.Unlock()

	p.logger.
		WithField("name", frontend.Name).
		WithField("url", frontend.URL).
		WithField("service_name", frontend.ServiceName).
		Debug("configuring frontend")

	if frontend.Certificate != nil {
		// configure HTTPS
		err := p.httpsServer.UpsertCertificate(
			frontend.URL.Host,
			frontend.Certificate)
		if err != nil {
			return err
		}

		err = p.httpsServer.UpsertRoute(
			frontend.URL,
			p.getServiceHandler(frontend.ServiceName))
		if err != nil {
			return err
		}

		// configure HTTP redirect
		httpURL := &url.URL{}
		*httpURL = *frontend.URL
		httpURL.Scheme = "http"

		err = p.httpServer.UpsertRoute(httpURL,
			http.RedirectHandler(
				frontend.URL.String(),
				http.StatusMovedPermanently))
		if err != nil {
			return err
		}
	} else {
		// configure HTTP
		err := p.httpServer.UpsertRoute(
			frontend.URL,
			p.getServiceHandler(frontend.ServiceName))
		if err != nil {
			return err
		}
	}

	return nil
}
