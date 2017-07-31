package startproxy

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"github.com/off-sync/platform-proxy-domain/frontends"
)

type dummyWebServer struct {
	FailAll bool
	routes  map[string]http.Handler
}

func (s *dummyWebServer) checkState() {
	if s.routes == nil {
		s.routes = make(map[string]http.Handler)
	}
}

func (s *dummyWebServer) UpsertRoute(route *url.URL, handler http.Handler) error {
	s.checkState()

	if s.FailAll {
		return fmt.Errorf("UpsertRoute(%v, %v)", route, handler)
	}

	s.routes[route.String()] = handler

	return nil
}

func (s *dummyWebServer) DeleteRoute(route *url.URL) {
	s.checkState()

	delete(s.routes, route.String())
}

func (s *dummyWebServer) UpsertCertificate(domainName string, cert *frontends.Certificate) error {
	if s.FailAll {
		return fmt.Errorf("UpsertCertificate(%s, %v)", domainName, cert)
	}

	return nil
}

type dummyResponseWriter struct {
	bytes.Buffer
	header http.Header
	status int
}

func newDummyResponseWriter() *dummyResponseWriter {
	return &dummyResponseWriter{
		header: make(map[string][]string),
	}
}

func (w *dummyResponseWriter) Header() http.Header {
	return w.header
}

func (w *dummyResponseWriter) WriteHeader(status int) {
	w.status = status
}

func (s *dummyWebServer) Handle(route *url.URL, r *http.Request) string {
	s.checkState()

	handler, found := s.routes[route.String()]
	if !found {
		return fmt.Sprintf("Not found: %s, got %v", route.String(), s.routes)
	}

	w := newDummyResponseWriter()
	handler.ServeHTTP(w, r)

	return string(w.Bytes())
}
