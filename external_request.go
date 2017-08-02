package slack

import "github.com/blendlabs/go-request"

// NewExternalRequest Creates a new external request
func NewExternalRequest() *request.Request {
	return request.New().WithMockProvider(request.MockedResponseInjector)
}
