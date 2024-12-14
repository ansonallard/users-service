package errors

import (
	"fmt"

	"github.com/ansonallard/users-service/internal/api"
)

type OAuth2Error struct {
	OAuth2Error api.OAuth2ErrorSchemaError `json:"error"`
	Err         error                      `json:"omitempty"`
}

func (e *OAuth2Error) Error() string {
	return fmt.Sprintf("%s", e.OAuth2Error)
}

// Implement Unwrap to support error wrapping (Go 1.13+)
func (e *OAuth2Error) Unwrap() error {
	return e.Err
}
