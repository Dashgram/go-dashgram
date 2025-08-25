package dashgram

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// TestSuite provides integration tests for the entire package
func TestSuite(t *testing.T) {
	t.Run("Integration", func(t *testing.T) {
		testIntegration(t)
	})

	t.Run("Concurrency", func(t *testing.T) {
		testConcurrency(t)
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		testErrorHandling(t)
	})
}

func testIntegration(t *testing.T) {
	helper := NewTestHelper()
	helper.AddResponse(200, `{"status":"success","details":"ok"}`)
	helper.AddResponse(200, `{"status":"success","details":"ok"}`)

	client := CreateTestClient(123, "test-key", WithHTTPClient(helper.MockHTTPClient()))
	defer client.Close()

	// Test synchronous operations
	err := client.TrackEvent(TestEventData)
	if err != nil {
		t.Errorf("TrackEvent failed: %v", err)
	}

	err = client.InvitedBy(TestUserData.UserID, TestUserData.InvitedBy)
	if err != nil {
		t.Errorf("InvitedBy failed: %v", err)
	}

	// Verify requests were made
	if helper.RequestCount != 2 {
		t.Errorf("expected 2 requests, got %d", helper.RequestCount)
	}
}

func testConcurrency(t *testing.T) {
	helper := NewTestHelper()

	// Add multiple responses for concurrent requests
	for i := 0; i < 10; i++ {
		helper.AddResponse(200, `{"status":"success","details":"ok"}`)
	}

	client := CreateTestClient(123, "test-key",
		WithHTTPClient(helper.MockHTTPClient()),
		WithUseAsync(),
		WithNumWorkers(3),
	)
	defer client.Close()

	// Send multiple concurrent requests
	for i := 0; i < 10; i++ {
		client.TrackEventAsync(map[string]any{
			"action": "concurrent_test",
			"index":  i,
		})
	}

	// Wait for all requests to be processed
	if !helper.WaitForRequests(10, 2*time.Second) {
		t.Errorf("not all requests were processed within timeout")
	}
}

func testErrorHandling(t *testing.T) {
	// Test forbidden error
	helper := NewTestHelper()
	helper.AddResponse(403, `{"status":"error","details":"forbidden"}`)

	client := CreateTestClient(123, "test-key", WithHTTPClient(helper.MockHTTPClient()))
	defer client.Close()

	err := client.TrackEvent(TestEventData)
	if err == nil {
		t.Errorf("expected error for forbidden response")
	}
	if !strings.Contains(err.Error(), "invalid credentials") {
		t.Errorf("expected error containing 'invalid credentials', got '%s'", err.Error())
	}

	// Test API error
	helper.Reset()
	helper.AddResponse(400, `{"status":"error","details":"bad request"}`)

	err = client.InvitedBy(TestUserData.UserID, TestUserData.InvitedBy)
	if err == nil {
		t.Errorf("expected error for bad request response")
	}
	if _, ok := err.(*DashgramAPIError); !ok {
		t.Errorf("expected DashgramAPIError, got %T", err)
	}

	// Test network error
	helper.Reset()
	helper.AddError(fmt.Errorf("network error"))

	err = client.TrackEvent(TestEventData)
	if err == nil {
		t.Errorf("expected error for network failure")
	}
	if !strings.Contains(err.Error(), "network error") {
		t.Errorf("expected network error message, got %s", err.Error())
	}
}

// Benchmark tests
func BenchmarkTrackEvent(b *testing.B) {
	helper := NewTestHelper()
	// Add responses for all benchmark iterations
	for i := 0; i < b.N; i++ {
		helper.AddResponse(200, `{"status":"success","details":"ok"}`)
	}

	client := CreateTestClient(123, "test-key", WithHTTPClient(helper.MockHTTPClient()))
	defer client.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := client.TrackEvent(TestEventData)
		if err != nil {
			b.Errorf("TrackEvent failed: %v", err)
		}
	}
}

func BenchmarkTrackEventAsync(b *testing.B) {
	helper := NewTestHelper()
	for i := 0; i < b.N; i++ {
		helper.AddResponse(200, `{"status":"success","details":"ok"}`)
	}

	client := CreateTestClient(123, "test-key",
		WithHTTPClient(helper.MockHTTPClient()),
		WithUseAsync(),
	)
	defer client.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.TrackEventAsync(TestEventData)
	}

	// Wait for all requests to complete
	if !helper.WaitForRequests(b.N, 5*time.Second) {
		b.Errorf("not all requests completed within timeout")
	}
}

func BenchmarkInvitedBy(b *testing.B) {
	helper := NewTestHelper()
	// Add responses for all benchmark iterations
	for i := 0; i < b.N; i++ {
		helper.AddResponse(200, `{"status":"success","details":"ok"}`)
	}

	client := CreateTestClient(123, "test-key", WithHTTPClient(helper.MockHTTPClient()))
	defer client.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := client.InvitedBy(TestUserData.UserID, TestUserData.InvitedBy)
		if err != nil {
			b.Errorf("InvitedBy failed: %v", err)
		}
	}
}

// Example tests
func ExampleDashgram_TrackEvent() {
	client := New(123, "your-access-key")
	defer client.Close()

	event := map[string]any{
		"action":  "page_view",
		"page":    "/home",
		"user_id": 12345,
	}

	err := client.TrackEvent(event)
	if err != nil {
		// Handle error
	}
}

func ExampleDashgram_TrackEventAsync() {
	client := New(123, "your-access-key", WithUseAsync())
	defer client.Close()

	event := map[string]any{
		"action":  "button_click",
		"button":  "signup",
		"user_id": 12345,
	}

	// This returns immediately, event is processed asynchronously
	client.TrackEventAsync(event)
}

func ExampleDashgram_InvitedBy() {
	client := New(123, "your-access-key")
	defer client.Close()

	err := client.InvitedBy(12345, 67890)
	if err != nil {
		// Handle error
	}
}
