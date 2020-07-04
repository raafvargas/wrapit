package http

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// MockedResponseResult ...
type MockedResponseResult struct {
	Body    string
	Status  int
	Headers http.Header
}

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// MockClient ...
type MockClient struct {
	client    *http.Client
	history   []string
	result    map[string]*MockedResponseResult
	roundTrip RoundTripFunc
}

// NewMockClient ...
func NewMockClient(options ...MockOption) *MockClient {
	mock := &MockClient{}
	mock.result = make(map[string]*MockedResponseResult)
	mock.roundTrip = mock.defaultRoundTrip

	for _, opt := range options {
		opt(mock)
	}

	mock.client = &http.Client{
		Transport: mock.roundTrip}

	return mock
}

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// Client ...
func (h *MockClient) Client() *http.Client {
	return h.client
}

// AddMock ...
func (h *MockClient) AddMock(uri string, result *MockedResponseResult) {
	WithMock(uri, result)(h)
}

// History ...
func (h *MockClient) History() []string {
	return h.history
}

func (h *MockClient) defaultRoundTrip(req *http.Request) *http.Response {
	uri := req.URL.RequestURI()

	h.history = append(h.history, uri)

	result, ok := h.result[uri]

	if !ok {
		result = &MockedResponseResult{}
		result.Body = "OK"
		result.Status = http.StatusOK
		result.Headers = make(http.Header)
	}

	return &http.Response{
		StatusCode: result.Status,
		Body:       ioutil.NopCloser(bytes.NewBufferString(result.Body)),
		Header:     result.Headers,
	}
}
