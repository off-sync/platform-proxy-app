package startproxy

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/off-sync/platform-proxy-domain/frontends"
)

type dummyWebServer struct {
	FailAll bool
	routes  map[string]http.Handler
}

func (s *dummyWebServer) UpsertRoute(route *url.URL, handler http.Handler) error {
	if s.FailAll {
		return fmt.Errorf("UpsertRoute(%v, %v)", route, handler)
	}

	return nil
}

func (s *dummyWebServer) DeleteRoute(route *url.URL) {
	// do nothing
}

func (s *dummyWebServer) UpsertCertificate(domainName string, cert *frontends.Certificate) error {
	if s.FailAll {
		return fmt.Errorf("UpsertCertificate(%s, %v)", domainName, cert)
	}

	return nil
}
