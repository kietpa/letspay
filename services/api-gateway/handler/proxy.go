package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// NewReverseProxy returns a reverse proxy handler
func NewReverseProxy(target string) http.HandlerFunc {
	parsedURL, err := url.Parse(target)
	if err != nil {
		panic("invalid proxy target: " + target)
	}

	proxy := httputil.NewSingleHostReverseProxy(parsedURL)

	proxy.ModifyResponse = func(resp *http.Response) error {
		if resp.Header.Get("Content-Type") == "" {
			resp.Header.Set("Content-Type", "application/json")
		}
		return nil
	}

	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Scheme = parsedURL.Scheme
		r.URL.Host = parsedURL.Host
		r.Host = parsedURL.Host
		r.URL.Path = singleJoiningSlash(parsedURL.Path, r.URL.Path)
		proxy.ServeHTTP(w, r)
	}
}

// Needed to handle joined paths properly
func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
