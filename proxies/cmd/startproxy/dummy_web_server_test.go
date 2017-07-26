package startproxy

import (
	"net/http"
	"net/url"

	"github.com/off-sync/platform-proxy-domain/frontends"
)

type dummyWebServer struct {
}

func (s *dummyWebServer) UpsertRoute(route *url.URL, handler http.Handler) error {
	// do nothing
	return nil
}

func (s *dummyWebServer) DeleteRoute(route *url.URL) {
	// do nothing
}

func (s *dummyWebServer) UpsertCertificate(domainName string, cert *frontends.Certificate) error {
	// do nothing
	return nil
}
