package dashgram

import (
	"testing"
)

func TestInvalidCredentialsError(t *testing.T) {
	err := &InvalidCredentialsError{}

	expected := "invalid credentials"
	if err.Error() != expected {
		t.Errorf("expected error message '%s', got '%s'", expected, err.Error())
	}
}

func TestDashgramAPIError(t *testing.T) {
	tests := []struct {
		name          string
		statusCode    int
		details       string
		expectedError string
	}{
		{
			name:          "basic API error",
			statusCode:    400,
			details:       "bad request",
			expectedError: "dashgram API error (status: 400): bad request",
		},
		{
			name:          "API error with empty details",
			statusCode:    500,
			details:       "",
			expectedError: "dashgram API error (status: 500): ",
		},
		{
			name:          "API error with special characters",
			statusCode:    403,
			details:       "forbidden: access denied",
			expectedError: "dashgram API error (status: 403): forbidden: access denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &DashgramAPIError{
				StatusCode: tt.statusCode,
				Details:    tt.details,
			}

			if err.Error() != tt.expectedError {
				t.Errorf("expected error message '%s', got '%s'", tt.expectedError, err.Error())
			}
		})
	}
}

func TestErrorTypeAssertions(t *testing.T) {
	// Test InvalidCredentialsError type assertion
	var err error = &InvalidCredentialsError{}

	if _, ok := err.(*InvalidCredentialsError); !ok {
		t.Errorf("failed to assert InvalidCredentialsError type")
	}

	// Test DashgramAPIError type assertion
	var apiErr error = &DashgramAPIError{
		StatusCode: 404,
		Details:    "not found",
	}

	if _, ok := apiErr.(*DashgramAPIError); !ok {
		t.Errorf("failed to assert DashgramAPIError type")
	}
}
