package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/mihailozarinschi/ports-app/internal/port"
)

var _ port.Repository = (*PortRepository)(nil)

// PortRepository is an in-memory implementation of the port.Repository
type PortRepository struct {
	rwm      sync.RWMutex
	portsMap map[string]port.Port
}

func NewPortRepository() port.Repository {
	return &PortRepository{
		rwm:      sync.RWMutex{},
		portsMap: make(map[string]port.Port),
	}
}

func (r *PortRepository) GetByID(ctx context.Context, portID string) (*port.Port, error) {
	r.rwm.RLock()
	defer r.rwm.RUnlock()

	p, found := r.portsMap[portID]
	if !found {
		return nil, fmt.Errorf("no port found with ID %q: %w", portID, port.ErrPortNotFound)
	}
	return &p, nil
}

func (r *PortRepository) CreateOrUpdate(ctx context.Context, port port.Port) (created bool, _ error) {
	r.rwm.Lock()
	defer r.rwm.Unlock()

	_, found := r.portsMap[port.ID]
	if !found {
		created = true
	}
	r.portsMap[port.ID] = port
	return created, nil
}
