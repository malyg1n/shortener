package errs

import "fmt"

// NotFoundError ...
type NotFoundError struct {
	StatusCode int
	Message    string
}

// Error return error message
func (e *NotFoundError) Error() string {
	return fmt.Sprintf("error: %s", e.Message)
}
