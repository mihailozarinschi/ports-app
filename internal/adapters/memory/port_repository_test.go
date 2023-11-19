package memory

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/mihailozarinschi/ports-app/internal/port"
)

// TestPortRepository starts an existing test suite that covers all
// expected behaviors from the port.Repository, and it injects the in-memory implementation.
func TestPortRepository(t *testing.T) {
	memRepo := NewPortRepository()
	suite.Run(t, &port.TestRepositorySuite{
		Repo: memRepo,
	})
}
