package handler

import (
	"net/http/httputil"
	"net/url"
)

func NewProxy(targetURL string) (*httputil.ReverseProxy, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	return httputil.NewSingleHostReverseProxy(target), nil
}
