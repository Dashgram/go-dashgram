package dashgram

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestDashgram_TrackEvent(t *testing.T) {
	tests := []struct {
		name          string
		event         any
		useAsync      bool
		mockResponse  *http.Response
		mockError     error
		expectedError string
		checkRequest  func(*http.Request) error
	}{
		{
			name:  "successful track event",
			event: map[string]string{"action": "click", "page": "home"},
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"status":"success","details":"ok"}`)),
			},
			checkRequest: func(req *http.Request) error {
				if !strings.HasSuffix(req.URL.Path, "/track") {
					return fmt.Errorf("expected endpoint '/track', got %s", req.URL.Path)
				}
				return nil
			},
		},
		{
			name:     "track event with async enabled",
			event:    map[string]string{"action": "view", "page": "about"},
			useAsync: true,
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"status":"success","details":"ok"}`)),
			},
		},
		{
			name:  "API error response",
			event: map[string]string{"action": "submit", "form": "contact"},
			mockResponse: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(strings.NewReader(`{"status":"error","details":"invalid event data"}`)),
			},
			expectedError: "dashgram API error (status: 400): invalid event data",
		},
		{
			name:          "network error",
			event:         map[string]string{"action": "load", "page": "dashboard"},
			mockError:     fmt.Errorf("connection timeout"),
			expectedError: "request failed: connection timeout",
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

			options := []Option{WithHTTPClient(mockClient)}
			if tt.useAsync {
				options = append(options, WithUseAsync())
			}

			d := New(123, "test-key", options...)
			defer d.Close()

			err := d.TrackEvent(tt.event)

			if tt.useAsync {
				// Async mode should return nil immediately
				if err != nil {
					t.Errorf("expected nil error for async mode, got %v", err)
				}
			} else if tt.expectedError != "" {
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

func TestDashgram_TrackEventWithContext(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		event         any
		useAsync      bool
		mockResponse  *http.Response
		mockError     error
		expectedError string
	}{
		{
			name:  "successful track event with context",
			ctx:   context.Background(),
			event: map[string]string{"action": "login", "method": "email"},
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"status":"success","details":"ok"}`)),
			},
		},
		{
			name:     "track event with context and async",
			ctx:      context.Background(),
			event:    map[string]string{"action": "logout"},
			useAsync: true,
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"status":"success","details":"ok"}`)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockHTTPClient{
				doFunc: func(req *http.Request) (*http.Response, error) {
					return tt.mockResponse, tt.mockError
				},
			}

			options := []Option{WithHTTPClient(mockClient)}
			if tt.useAsync {
				options = append(options, WithUseAsync())
			}

			d := New(123, "test-key", options...)
			defer d.Close()

			err := d.TrackEventWithContext(tt.ctx, tt.event)

			if tt.useAsync {
				if err != nil {
					t.Errorf("expected nil error for async mode, got %v", err)
				}
			} else if tt.expectedError != "" {
				if err == nil {
					t.Errorf("expected error '%s', got nil", tt.expectedError)
				} else if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("expected error containing '%s', got '%s'", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestDashgram_InvitedBy(t *testing.T) {
	tests := []struct {
		name          string
		userID        int
		invitedBy     int
		useAsync      bool
		mockResponse  *http.Response
		mockError     error
		expectedError string
		checkRequest  func(*http.Request) error
	}{
		{
			name:      "successful invited by",
			userID:    12345,
			invitedBy: 67890,
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"status":"success","details":"ok"}`)),
			},
			checkRequest: func(req *http.Request) error {
				if !strings.HasSuffix(req.URL.Path, "/invited_by") {
					return fmt.Errorf("expected endpoint '/invited_by', got %s", req.URL.Path)
				}
				return nil
			},
		},
		{
			name:      "invited by with async enabled",
			userID:    11111,
			invitedBy: 22222,
			useAsync:  true,
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"status":"success","details":"ok"}`)),
			},
		},
		{
			name:      "API error response",
			userID:    33333,
			invitedBy: 44444,
			mockResponse: &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(strings.NewReader(`{"status":"error","details":"user not found"}`)),
			},
			expectedError: "dashgram API error (status: 404): user not found",
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

			options := []Option{WithHTTPClient(mockClient)}
			if tt.useAsync {
				options = append(options, WithUseAsync())
			}

			d := New(123, "test-key", options...)
			defer d.Close()

			err := d.InvitedBy(tt.userID, tt.invitedBy)

			if tt.useAsync {
				if err != nil {
					t.Errorf("expected nil error for async mode, got %v", err)
				}
			} else if tt.expectedError != "" {
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

func TestDashgram_InvitedByWithContext(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		userID        int
		invitedBy     int
		useAsync      bool
		mockResponse  *http.Response
		mockError     error
		expectedError string
	}{
		{
			name:      "successful invited by with context",
			ctx:       context.Background(),
			userID:    55555,
			invitedBy: 66666,
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"status":"success","details":"ok"}`)),
			},
		},
		{
			name:      "invited by with context and async",
			ctx:       context.Background(),
			userID:    77777,
			invitedBy: 88888,
			useAsync:  true,
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"status":"success","details":"ok"}`)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockHTTPClient{
				doFunc: func(req *http.Request) (*http.Response, error) {
					return tt.mockResponse, tt.mockError
				},
			}

			options := []Option{WithHTTPClient(mockClient)}
			if tt.useAsync {
				options = append(options, WithUseAsync())
			}

			d := New(123, "test-key", options...)
			defer d.Close()

			err := d.InvitedByWithContext(tt.ctx, tt.userID, tt.invitedBy)

			if tt.useAsync {
				if err != nil {
					t.Errorf("expected nil error for async mode, got %v", err)
				}
			} else if tt.expectedError != "" {
				if err == nil {
					t.Errorf("expected error '%s', got nil", tt.expectedError)
				} else if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("expected error containing '%s', got '%s'", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
