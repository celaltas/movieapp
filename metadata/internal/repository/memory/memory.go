package memory

import (
	"context"
	"sync"

	"movieexample.com/metadata/internal/repository"
	model "movieexample.com/metadata/pkg/model"
)

// Repository defines a memory movie metadata repository.
type Repository struct {
	sync.RWMutex
	data map[string]*model.Metadata
}

// New creates a new memory repository.
func New() *Repository {
	return &Repository{
		data: make(map[string]*model.Metadata),
	}
}

// Get retrieves movie metadata for by movie id.
func (r *Repository) Get(_ context.Context, id string) (*model.Metadata, error) {
	r.RLock()
	defer r.RUnlock()
	if val, ok := r.data[id]; ok {
		return val, nil
	}
	return nil, repository.ErrNotFound
}

// Put adds movie metadata for a given movie id.
func (r *Repository) Put(_ context.Context, m *model.Metadata) error {
	r.Lock()
	defer r.Unlock()
	r.data[m.ID] = m
	return nil
}
