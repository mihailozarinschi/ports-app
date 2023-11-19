package port

import (
	"context"
	"errors"
)

var (
	ErrPortNotFound = errors.New("port not found")
	// ErrConnectionFailed is returned when using an external storage, and connection fails during requests to it.
	ErrConnectionFailed = errors.New("connection to external service failed")
)

type Repository interface {
	// GetByID returns a Port given its ID.
	// ErrPortNotFound is returned if nothing was found.
	GetByID(ctx context.Context, portID string) (*Port, error)
	// CreateOrUpdate will attempt to create a new entry for the given port,
	// or update it if already exists.
	CreateOrUpdate(ctx context.Context, port Port) (created bool, _ error)
	// Implement Delete and other queries when needed
}
