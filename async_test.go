package dashgram

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestDashgram_TrackEventAsync(t *testing.T) {
	tests := []struct {
		name          string
		event         any
		mockResponse  *http.Response
		mockError     error
		expectedError string
	}{
		{
			name:  "successful async track event",
			event: map[string]string{"action": "async_click", "page": "async_home"},
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"status":"success","details":"ok"}`)),
			},
		},
		{
			name:  "async track event with API error",
			event: map[string]string{"action": "async_submit", "form": "async_contact"},
			mockResponse: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(strings.NewReader(`{"status":"error","details":"invalid async event data"}`)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requestCount int
			var mu sync.Mutex

			mockClient := &mockHTTPClient{
				doFunc: func(req *http.Request) (*http.Response, error) {
					mu.Lock()
					requestCount++
					mu.Unlock()
					return tt.mockResponse, tt.mockError
				},
			}

			d := New(123, "test-key", WithHTTPClient(mockClient), WithUseAsync())
			defer d.Close()

			// Enqueue the async task
			d.TrackEventAsync(tt.event)

			// Give some time for the worker to process the task
			time.Sleep(50 * time.Millisecond)

			// Check that the request was made
			mu.Lock()
			if requestCount == 0 {
				t.Errorf("expected request to be made, but none was")
			}
			mu.Unlock()
		})
	}
}

func TestDashgram_TrackEventAsyncWithContext(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		event         any
		mockResponse  *http.Response
		mockError     error
		expectedError string
	}{
		{
			name:  "successful async track event with context",
			ctx:   context.Background(),
			event: map[string]string{"action": "async_login", "method": "async_email"},
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"status":"success","details":"ok"}`)),
			},
		},
		{
			name:  "async track event with cancelled context",
			ctx:   func() context.Context { ctx, cancel := context.WithCancel(context.Background()); cancel(); return ctx }(),
			event: map[string]string{"action": "async_test"},
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"status":"success","details":"ok"}`)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requestCount int
			var mu sync.Mutex

			mockClient := &mockHTTPClient{
				doFunc: func(req *http.Request) (*http.Response, error) {
					mu.Lock()
					requestCount++
					mu.Unlock()
					return tt.mockResponse, tt.mockError
				},
			}

			d := New(123, "test-key", WithHTTPClient(mockClient), WithUseAsync())
			defer d.Close()

			// Enqueue the async task with context
			d.TrackEventAsyncWithContext(tt.ctx, tt.event)

			// Give some time for the worker to process the task
			time.Sleep(50 * time.Millisecond)

			// For cancelled context, we might not see a request
			if tt.ctx.Err() == nil {
				mu.Lock()
				if requestCount == 0 {
					t.Errorf("expected request to be made, but none was")
				}
				mu.Unlock()
			}
		})
	}
}

func TestDashgram_InvitedByAsync(t *testing.T) {
	tests := []struct {
		name          string
		userID        int
		invitedBy     int
		mockResponse  *http.Response
		mockError     error
		expectedError string
	}{
		{
			name:      "successful async invited by",
			userID:    12345,
			invitedBy: 67890,
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"status":"success","details":"ok"}`)),
			},
		},
		{
			name:      "async invited by with API error",
			userID:    33333,
			invitedBy: 44444,
			mockResponse: &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(strings.NewReader(`{"status":"error","details":"async user not found"}`)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requestCount int
			var mu sync.Mutex

			mockClient := &mockHTTPClient{
				doFunc: func(req *http.Request) (*http.Response, error) {
					mu.Lock()
					requestCount++
					mu.Unlock()
					return tt.mockResponse, tt.mockError
				},
			}

			d := New(123, "test-key", WithHTTPClient(mockClient), WithUseAsync())
			defer d.Close()

			// Enqueue the async task
			d.InvitedByAsync(tt.userID, tt.invitedBy)

			// Give some time for the worker to process the task
			time.Sleep(50 * time.Millisecond)

			// Check that the request was made
			mu.Lock()
			if requestCount == 0 {
				t.Errorf("expected request to be made, but none was")
			}
			mu.Unlock()
		})
	}
}

func TestDashgram_InvitedByAsyncWithContext(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		userID        int
		invitedBy     int
		mockResponse  *http.Response
		mockError     error
		expectedError string
	}{
		{
			name:      "successful async invited by with context",
			ctx:       context.Background(),
			userID:    55555,
			invitedBy: 66666,
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"status":"success","details":"ok"}`)),
			},
		},
		{
			name:      "async invited by with cancelled context",
			ctx:       func() context.Context { ctx, cancel := context.WithCancel(context.Background()); cancel(); return ctx }(),
			userID:    99999,
			invitedBy: 10000,
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"status":"success","details":"ok"}`)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requestCount int
			var mu sync.Mutex

			mockClient := &mockHTTPClient{
				doFunc: func(req *http.Request) (*http.Response, error) {
					mu.Lock()
					requestCount++
					mu.Unlock()
					return tt.mockResponse, tt.mockError
				},
			}

			d := New(123, "test-key", WithHTTPClient(mockClient), WithUseAsync())
			defer d.Close()

			// Enqueue the async task with context
			d.InvitedByAsyncWithContext(tt.ctx, tt.userID, tt.invitedBy)

			// Give some time for the worker to process the task
			time.Sleep(50 * time.Millisecond)

			// For cancelled context, we might not see a request
			if tt.ctx.Err() == nil {
				mu.Lock()
				if requestCount == 0 {
					t.Errorf("expected request to be made, but none was")
				}
				mu.Unlock()
			}
		})
	}
}

func TestDashgram_enqueueTask(t *testing.T) {
	d := New(123, "test-key", WithUseAsync())
	defer d.Close()

	// Test enqueueing a task
	task := asyncTask{
		ctx:      context.Background(),
		endpoint: "test",
		data:     map[string]string{"test": "data"},
	}

	// This should not block
	d.enqueueTask(task)

	// Give some time for the task to be processed
	time.Sleep(10 * time.Millisecond)
}

func TestDashgram_AsyncWorkerShutdown(t *testing.T) {
	d := New(123, "test-key", WithUseAsync())

	// Enqueue a task
	d.TrackEventAsync(map[string]string{"action": "test"})

	// Close the client
	d.Close()

	// Try to enqueue another task after shutdown
	d.TrackEventAsync(map[string]string{"action": "test_after_shutdown"})

	// This should not panic or block
	time.Sleep(10 * time.Millisecond)
}

func TestDashgram_MultipleAsyncWorkers(t *testing.T) {
	var requestCount int
	var mu sync.Mutex

	mockClient := &mockHTTPClient{
		doFunc: func(req *http.Request) (*http.Response, error) {
			mu.Lock()
			requestCount++
			mu.Unlock()
			// Simulate some processing time
			time.Sleep(10 * time.Millisecond)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"status":"success","details":"ok"}`)),
			}, nil
		},
	}

	d := New(123, "test-key", WithHTTPClient(mockClient), WithUseAsync(), WithNumWorkers(3))
	defer d.Close()

	// Enqueue multiple tasks
	for i := 0; i < 5; i++ {
		d.TrackEventAsync(map[string]string{"action": "test", "index": fmt.Sprintf("%d", i)})
	}

	// Give time for all tasks to be processed
	time.Sleep(200 * time.Millisecond)

	// Check that all requests were made
	mu.Lock()
	if requestCount != 5 {
		t.Errorf("expected 5 requests, got %d", requestCount)
	}
	mu.Unlock()
}
