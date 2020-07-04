package http

import "net/http"

// MockOption ...
type MockOption func(*MockClient)

// WithMock ...
func WithMock(uri string, result *MockedResponseResult) MockOption {
	return func(mock *MockClient) {
		if result.Headers == nil {
			result.Headers = make(http.Header)
		}

		mock.result[uri] = result
	}
}

// WithRoundTrip ...
func WithRoundTrip(roundTrip RoundTripFunc) MockOption {
	return func(mock *MockClient) {
		mock.roundTrip = roundTrip
	}
}
