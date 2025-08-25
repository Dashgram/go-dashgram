package dashgram

import (
	"encoding/json"
	"testing"
)

func TestTrackEventRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  TrackEventRequest
		expected string
	}{
		{
			name: "basic track event request",
			request: TrackEventRequest{
				Updates: []any{
					map[string]string{"action": "click", "page": "home"},
				},
				Origin: "Test App",
			},
			expected: `{"updates":[{"action":"click","page":"home"}],"origin":"Test App"}`,
		},
		{
			name: "track event request without origin",
			request: TrackEventRequest{
				Updates: []any{
					map[string]string{"action": "view", "page": "about"},
					map[string]string{"action": "submit", "form": "contact"},
				},
			},
			expected: `{"updates":[{"action":"view","page":"about"},{"action":"submit","form":"contact"}]}`,
		},
		{
			name: "track event request with complex data",
			request: TrackEventRequest{
				Updates: []any{
					map[string]any{
						"action":   "purchase",
						"amount":   99.99,
						"currency": "USD",
						"items":    []string{"item1", "item2"},
					},
				},
				Origin: "E-commerce App",
			},
			expected: `{"updates":[{"action":"purchase","amount":99.99,"currency":"USD","items":["item1","item2"]}],"origin":"E-commerce App"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.request)
			if err != nil {
				t.Errorf("failed to marshal TrackEventRequest: %v", err)
			}

			if string(data) != tt.expected {
				t.Errorf("expected JSON '%s', got '%s'", tt.expected, string(data))
			}

			// Test unmarshaling back
			var unmarshaled TrackEventRequest
			if err := json.Unmarshal(data, &unmarshaled); err != nil {
				t.Errorf("failed to unmarshal TrackEventRequest: %v", err)
			}
		})
	}
}

func TestInvitedByRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  InvitedByRequest
		expected string
	}{
		{
			name: "basic invited by request",
			request: InvitedByRequest{
				UserID:    12345,
				InvitedBy: 67890,
				Origin:    "Test App",
			},
			expected: `{"user_id":12345,"invited_by":67890,"origin":"Test App"}`,
		},
		{
			name: "invited by request without origin",
			request: InvitedByRequest{
				UserID:    11111,
				InvitedBy: 22222,
			},
			expected: `{"user_id":11111,"invited_by":22222}`,
		},
		{
			name: "invited by request with large IDs",
			request: InvitedByRequest{
				UserID:    999999999,
				InvitedBy: 888888888,
				Origin:    "Large Scale App",
			},
			expected: `{"user_id":999999999,"invited_by":888888888,"origin":"Large Scale App"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.request)
			if err != nil {
				t.Errorf("failed to marshal InvitedByRequest: %v", err)
			}

			if string(data) != tt.expected {
				t.Errorf("expected JSON '%s', got '%s'", tt.expected, string(data))
			}

			// Test unmarshaling back
			var unmarshaled InvitedByRequest
			if err := json.Unmarshal(data, &unmarshaled); err != nil {
				t.Errorf("failed to unmarshal InvitedByRequest: %v", err)
			}

			// Verify the unmarshaled data matches the original
			if unmarshaled.UserID != tt.request.UserID {
				t.Errorf("expected UserID %d, got %d", tt.request.UserID, unmarshaled.UserID)
			}
			if unmarshaled.InvitedBy != tt.request.InvitedBy {
				t.Errorf("expected InvitedBy %d, got %d", tt.request.InvitedBy, unmarshaled.InvitedBy)
			}
			if unmarshaled.Origin != tt.request.Origin {
				t.Errorf("expected Origin '%s', got '%s'", tt.request.Origin, unmarshaled.Origin)
			}
		})
	}
}

func TestRequestStructTags(t *testing.T) {
	// Test that the JSON tags are working correctly
	trackRequest := TrackEventRequest{
		Updates: []any{map[string]string{"test": "data"}},
		Origin:  "Test Origin",
	}

	data, err := json.Marshal(trackRequest)
	if err != nil {
		t.Errorf("failed to marshal: %v", err)
	}

	// Check that the JSON contains the expected field names
	jsonStr := string(data)
	if !contains(jsonStr, "updates") {
		t.Errorf("expected JSON to contain 'updates' field")
	}
	if !contains(jsonStr, "origin") {
		t.Errorf("expected JSON to contain 'origin' field")
	}

	invitedByRequest := InvitedByRequest{
		UserID:    123,
		InvitedBy: 456,
		Origin:    "Test Origin",
	}

	data, err = json.Marshal(invitedByRequest)
	if err != nil {
		t.Errorf("failed to marshal: %v", err)
	}

	jsonStr = string(data)
	if !contains(jsonStr, "user_id") {
		t.Errorf("expected JSON to contain 'user_id' field")
	}
	if !contains(jsonStr, "invited_by") {
		t.Errorf("expected JSON to contain 'invited_by' field")
	}
	if !contains(jsonStr, "origin") {
		t.Errorf("expected JSON to contain 'origin' field")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			contains(s[1:len(s)-1], substr)))
}
