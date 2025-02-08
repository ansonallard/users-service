package errors

import (
	"fmt"

	"github.com/ansonallard/users-service/internal/api"
)

type OAuth2Error struct {
	OAuth2Error api.OAuth2ErrorSchemaError `json:"error"`
	Err         error                      `json:",omitempty"`
}

func (e *OAuth2Error) Error() string {
	return fmt.Sprintf("%s", e.OAuth2Error)
}

// Implement Unwrap to support error wrapping (Go 1.13+)
func (e *OAuth2Error) Unwrap() error {
	return e.Err
}

type UserNotFoundError struct {
	Usr      string
	TenantId string
	Err      error
}

func (e UserNotFoundError) Error() string {
	return fmt.Sprintf("user '%s' not found in tenant '%s'", e.Usr, e.TenantId)
}

type TenantNotFoundError struct {
	Id string
}

func (t TenantNotFoundError) Error() string {
	return fmt.Sprintf("tenant not found with id %s", t.Id)
}

type UserExistsError struct {
}

func (u UserExistsError) Error() string {
	return "error: user already exists"
}

type NotAuthorizedError struct{}

func (n NotAuthorizedError) Error() string {
	return "Not Authorized"
}
