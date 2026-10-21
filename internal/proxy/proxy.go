package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Proxy interface {
	NewReverseProxy(target string) (*httputil.ReverseProxy, error)
	ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request)
}

type proxy struct{}

func NewProxy() Proxy {
	return &proxy{}
}

func (pxy *proxy) NewReverseProxy(target string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	rp := httputil.NewSingleHostReverseProxy(url)
	rp.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		http.Error(w, "backend unavailable", http.StatusBadGateway)
	}

	return rp, nil
}

func (pxy *proxy) ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}
