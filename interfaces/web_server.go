package interfaces

import (
	"net/http"
	"net/url"
)

// WebServer defines an interface provides methods to upsert and delete routes
// to a handler.
type WebServer interface {
	// UpsertRoute adds a route to the web server, forwarding all requests to the
	// provided handler. It returns an error if either parameter is nil.
	UpsertRoute(route *url.URL, handler http.Handler) error

	// DeleteRoute deletes a route from the web server.
	DeleteRoute(route *url.URL)
}
