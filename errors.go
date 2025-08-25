package dashgram

import "fmt"

// InvalidCredentialsError represents an invalid credentials error
type InvalidCredentialsError struct{}

func (e *InvalidCredentialsError) Error() string {
	return "invalid credentials"
}

// DashgramAPIError represents an API error from Dashgram
type DashgramAPIError struct {
	StatusCode int
	Details    string
}

func (e *DashgramAPIError) Error() string {
	return fmt.Sprintf("dashgram API error (status: %d): %s", e.StatusCode, e.Details)
}
