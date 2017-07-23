package startproxy

import (
	"net/http"
	"net/url"
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
