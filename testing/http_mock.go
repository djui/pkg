package testing

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// HTTPMock keeps track of registered routes and the latest request sent.
type HTTPMock struct {
	LatestRequest *http.Request

	registrar map[route]RoundTripFunc
}

// RoundTripFunc is a convenience type from the http.RoundTripper interface. It
// describes a usual HTTP request-response round-trip.
type RoundTripFunc func(*http.Request) (*http.Response, error)

type route struct {
	method string
	path   string
}

// Register adds (or overrides) a mocked response to a given route. A route
// consists of an HTTP method and an URL root path.
func (m *HTTPMock) Register(method, path string, responseFunc RoundTripFunc) {
	if m.registrar == nil {
		m.registrar = make(map[route]RoundTripFunc)
	}

	m.registrar[route{method, path}] = responseFunc
}

// RoundTrip satisfies the http.RoundTripper interface.
func (m *HTTPMock) RoundTrip(req *http.Request) (*http.Response, error) {
	m.LatestRequest = req

	reqRoute := route{req.Method, req.URL.Path}
	respFunc, ok := m.registrar[reqRoute]
	if !ok {
		return FakeResponse(req, http.StatusNotFound, nil)
	}

	return respFunc(req)
}

// FakeResponse constructs a simple mocked response given a status code and
// body. The content length is infered from the body length.
func FakeResponse(req *http.Request, statusCode int, body []byte) (*http.Response, error) {
	return &http.Response{
		StatusCode:    statusCode,
		Status:        http.StatusText(statusCode),
		Request:       req,
		Body:          ioutil.NopCloser(bytes.NewReader(body)),
		ContentLength: int64(len(body)),
	}, nil
}
