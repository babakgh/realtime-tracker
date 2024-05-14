package main

type Request interface {
	// Get Key
}

type MockRequest struct {
	Key       string
	BackendID int
}

func NewMockRequest(key string, backendID int) *MockRequest {
	return &MockRequest{
		Key:       key,
		BackendID: backendID,
	}
}

// {a, b, c, d, e, f, g, h, i, j}
// {0, 1, 0, 1, 0, 1, 0, 1, 0, 1}
