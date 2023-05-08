package http

import (
	"net/http"
	"time"
)

type RequestOption func(*Request)

func WithHost(host string) RequestOption {
	return func(p *Request) {
		p.Host = host
	}
}

func WithUrl(url string) RequestOption {
	return func(p *Request) {
		p.Url = url
	}
}

func WithBodySize(bodySize int) RequestOption {
	return func(p *Request) {
		p.bodySize = bodySize
	}
}

func WithMethod(method string) RequestOption {
	return func(p *Request) {
		p.Method = method
	}
}

func WithTransport(transport *http.Transport) RequestOption {
	return func(p *Request) {
		p.Transport = transport
	}
}

func WithTimeout(timeout time.Duration) RequestOption {
	return func(p *Request) {
		p.Timeout = timeout
	}
}

func WithRequestType(requestType RequestType) RequestOption {
	return func(p *Request) {
		if _, ok := types[requestType]; ok {
			p.requestType = requestType
		}
	}
}
