package dashgram

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name       string
		projectID  int
		accessKey  string
		options    []Option
		wantErr    bool
		checkFuncs []func(*Dashgram) error
	}{
		{
			name:      "basic initialization",
			projectID: 123,
			accessKey: "test-key",
			checkFuncs: []func(*Dashgram) error{
				func(d *Dashgram) error {
					if d.ProjectID != 123 {
						t.Errorf("expected ProjectID 123, got %d", d.ProjectID)
					}
					if d.AccessKey != "test-key" {
						t.Errorf("expected AccessKey 'test-key', got %s", d.AccessKey)
					}
					if d.APIURL != "https://api.dashgram.io/v1/123" {
						t.Errorf("expected APIURL 'https://api.dashgram.io/v1/123', got %s", d.APIURL)
					}
					if d.Origin != "Go + Dashgram SDK" {
						t.Errorf("expected Origin 'Go + Dashgram SDK', got %s", d.Origin)
					}
					if d.useAsync != false {
						t.Errorf("expected useAsync false, got %v", d.useAsync)
					}
					if d.numWorkers != 1 {
						t.Errorf("expected numWorkers 1, got %d", d.numWorkers)
					}
					return nil
				},
			},
		},
		{
			name:      "with custom API URL",
			projectID: 456,
			accessKey: "test-key-2",
			options: []Option{
				WithAPIURL("https://custom.api.com/v2"),
			},
			checkFuncs: []func(*Dashgram) error{
				func(d *Dashgram) error {
					if d.APIURL != "https://custom.api.com/v2/456" {
						t.Errorf("expected APIURL 'https://custom.api.com/v2/456', got %s", d.APIURL)
					}
					return nil
				},
			},
		},
		{
			name:      "with custom origin",
			projectID: 789,
			accessKey: "test-key-3",
			options: []Option{
				WithOrigin("Custom App"),
			},
			checkFuncs: []func(*Dashgram) error{
				func(d *Dashgram) error {
					if d.Origin != "Custom App" {
						t.Errorf("expected Origin 'Custom App', got %s", d.Origin)
					}
					return nil
				},
			},
		},
		{
			name:      "with async enabled",
			projectID: 999,
			accessKey: "test-key-4",
			options: []Option{
				WithUseAsync(),
			},
			checkFuncs: []func(*Dashgram) error{
				func(d *Dashgram) error {
					if d.useAsync != true {
						t.Errorf("expected useAsync true, got %v", d.useAsync)
					}
					return nil
				},
			},
		},
		{
			name:      "with multiple workers",
			projectID: 111,
			accessKey: "test-key-5",
			options: []Option{
				WithNumWorkers(5),
			},
			checkFuncs: []func(*Dashgram) error{
				func(d *Dashgram) error {
					if d.numWorkers != 5 {
						t.Errorf("expected numWorkers 5, got %d", d.numWorkers)
					}
					return nil
				},
			},
		},
		{
			name:      "with custom HTTP client",
			projectID: 222,
			accessKey: "test-key-6",
			options: []Option{
				WithHTTPClient(&http.Client{Timeout: 60 * time.Second}),
			},
			checkFuncs: []func(*Dashgram) error{
				func(d *Dashgram) error {
					if d.client == nil {
						t.Errorf("expected HTTP client to be set")
					}
					return nil
				},
			},
		},
		{
			name:      "with all options",
			projectID: 333,
			accessKey: "test-key-7",
			options: []Option{
				WithAPIURL("https://test.api.com/v3"),
				WithOrigin("Test App"),
				WithUseAsync(),
				WithNumWorkers(3),
				WithHTTPClient(&http.Client{Timeout: 45 * time.Second}),
			},
			checkFuncs: []func(*Dashgram) error{
				func(d *Dashgram) error {
					if d.APIURL != "https://test.api.com/v3/333" {
						t.Errorf("expected APIURL 'https://test.api.com/v3/333', got %s", d.APIURL)
					}
					if d.Origin != "Test App" {
						t.Errorf("expected Origin 'Test App', got %s", d.Origin)
					}
					if d.useAsync != true {
						t.Errorf("expected useAsync true, got %v", d.useAsync)
					}
					if d.numWorkers != 3 {
						t.Errorf("expected numWorkers 3, got %d", d.numWorkers)
					}
					return nil
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := New(tt.projectID, tt.accessKey, tt.options...)

			// Run all check functions
			for _, checkFunc := range tt.checkFuncs {
				if err := checkFunc(d); err != nil {
					t.Errorf("check failed: %v", err)
				}
			}

			// Clean up
			d.Close()
		})
	}
}

func TestDashgram_Close(t *testing.T) {
	d := New(123, "test-key", WithUseAsync())

	// Verify worker is running
	if d.workerCtx.Err() != nil {
		t.Errorf("expected worker context to be active")
	}

	// Close the client
	d.Close()

	// Verify worker is stopped
	if d.workerCtx.Err() == nil {
		t.Errorf("expected worker context to be cancelled after Close()")
	}
}

func TestDashgram_StartWorker(t *testing.T) {
	d := New(123, "test-key", WithUseAsync())
	defer d.Close()

	// Verify worker is running
	if d.workerCtx.Err() != nil {
		t.Errorf("expected worker context to be active after StartWorker()")
	}

	// Test that worker can process tasks
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Send a task to the worker
	d.enqueueTask(asyncTask{
		ctx:      ctx,
		endpoint: "test",
		data:     map[string]string{"test": "data"},
	})

	// Give some time for the task to be processed
	time.Sleep(10 * time.Millisecond)
}

func TestDashgram_request(t *testing.T) {
	tests := []struct {
		name          string
		endpoint      string
		data          any
		mockResponse  *http.Response
		mockError     error
		expectedError string
		checkRequest  func(*http.Request) error
	}{
		{
			name:     "successful request",
			endpoint: "track",
			data:     map[string]string{"event": "test"},
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"status":"success","details":"ok"}`)),
			},
			checkRequest: func(req *http.Request) error {
				if req.Method != "POST" {
					return fmt.Errorf("expected POST method, got %s", req.Method)
				}
				if req.Header.Get("Authorization") != "Bearer test-key" {
					return fmt.Errorf("expected Authorization header 'Bearer test-key', got %s", req.Header.Get("Authorization"))
				}
				if req.Header.Get("Content-Type") != "application/json" {
					return fmt.Errorf("expected Content-Type header 'application/json', got %s", req.Header.Get("Content-Type"))
				}
				return nil
			},
		},
		{
			name:     "forbidden response",
			endpoint: "track",
			data:     map[string]string{"event": "test"},
			mockResponse: &http.Response{
				StatusCode: http.StatusForbidden,
				Body:       io.NopCloser(strings.NewReader(`{"status":"error","details":"forbidden"}`)),
			},
			expectedError: "invalid credentials",
		},
		{
			name:     "API error response",
			endpoint: "track",
			data:     map[string]string{"event": "test"},
			mockResponse: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(strings.NewReader(`{"status":"error","details":"bad request"}`)),
			},
			expectedError: "dashgram API error (status: 400): bad request",
		},
		{
			name:          "network error",
			endpoint:      "track",
			data:          map[string]string{"event": "test"},
			mockError:     fmt.Errorf("network error"),
			expectedError: "request failed: network error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockHTTPClient{
				doFunc: func(req *http.Request) (*http.Response, error) {
					if tt.checkRequest != nil {
						if err := tt.checkRequest(req); err != nil {
							t.Errorf("request check failed: %v", err)
						}
					}
					return tt.mockResponse, tt.mockError
				},
			}

			d := New(123, "test-key", WithHTTPClient(mockClient))
			defer d.Close()

			err := d.request(context.Background(), tt.endpoint, tt.data)

			if tt.expectedError != "" {
				if err == nil {
					t.Errorf("expected error '%s', got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError {
					t.Errorf("expected error '%s', got '%s'", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
