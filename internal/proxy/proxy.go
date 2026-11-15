package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewProxy(server string) (http.Handler, error) {
	targetURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	return httputil.NewSingleHostReverseProxy(targetURL), nil
}
