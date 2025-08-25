package dashgram

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Mock HTTP client for testing
type mockHTTPClient struct {
	doFunc func(*http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.doFunc(req)
}

// TestHelper provides common test utilities
type TestHelper struct {
	RequestCount int
	mu           sync.Mutex
	Responses    []*http.Response
	Errors       []error
}

// NewTestHelper creates a new test helper instance
func NewTestHelper() *TestHelper {
	return &TestHelper{
		Responses: make([]*http.Response, 0),
		Errors:    make([]error, 0),
	}
}

// MockHTTPClient creates a mock HTTP client for testing
func (th *TestHelper) MockHTTPClient() *mockHTTPClient {
	return &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			th.mu.Lock()
			th.RequestCount++
			responseIndex := th.RequestCount - 1
			th.mu.Unlock()

			var response *http.Response
			var err error

			if responseIndex < len(th.Responses) {
				response = th.Responses[responseIndex]
			}
			if responseIndex < len(th.Errors) {
				err = th.Errors[responseIndex]
			}

			return response, err
		},
	}
}

// AddResponse adds a response to the mock client
func (th *TestHelper) AddResponse(statusCode int, body string) {
	th.Responses = append(th.Responses, &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
	})
}

// AddError adds an error to the mock client
func (th *TestHelper) AddError(err error) {
	th.Errors = append(th.Errors, err)
}

// Reset resets the test helper state
func (th *TestHelper) Reset() {
	th.mu.Lock()
	defer th.mu.Unlock()
	th.RequestCount = 0
	th.Responses = make([]*http.Response, 0)
	th.Errors = make([]error, 0)
}

// WaitForRequests waits for a specified number of requests to be made
func (th *TestHelper) WaitForRequests(expected int, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		th.mu.Lock()
		count := th.RequestCount
		th.mu.Unlock()
		if count >= expected {
			return true
		}
		time.Sleep(10 * time.Millisecond)
	}
	return false
}

// CreateTestClient creates a Dashgram client with test configuration
func CreateTestClient(projectID int, accessKey string, options ...Option) *Dashgram {
	// Set default test options if none provided
	if len(options) == 0 {
		options = append(options, WithHTTPClient(&http.Client{Timeout: 5 * time.Second}))
	}

	return New(projectID, accessKey, options...)
}

// CreateTestContext creates a test context with timeout
func CreateTestContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

// TestEventData provides common test event data
var TestEventData = map[string]any{
	"action":    "test_action",
	"page":      "test_page",
	"user_id":   12345,
	"timestamp": time.Now().Unix(),
}

// TestUserData provides common test user data
var TestUserData = struct {
	UserID    int
	InvitedBy int
}{
	UserID:    12345,
	InvitedBy: 67890,
}
