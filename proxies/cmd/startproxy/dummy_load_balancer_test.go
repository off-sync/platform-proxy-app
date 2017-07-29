package startproxy

import (
	"fmt"
	"net/http"
	"net/url"
)

type dummyLoadBalancer struct{}

func (lb *dummyLoadBalancer) UpsertService(name string, urls ...*url.URL) (http.Handler, error) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Service: %s\n\tURLs: %v", name, urls)
	}), nil
}
