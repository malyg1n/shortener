package errs

import "fmt"

// ValidationError ...
type ValidationError struct {
	StatusCode int
	Message    string
}

// Error return error message
func (e *ValidationError) Error() string {
	return fmt.Sprintf("error: %s", e.Message)
}
