package metadata

import (
	"context"
	"errors"

	"movieexample.com/metadata/internal/repository"
	model "movieexample.com/metadata/pkg/model"
)

// ErrNotFound is returned when a requested record is not found.
var ErrNotFound = errors.New("not found")

type metadataRepository interface {
	Get(context.Context, string) (*model.Metadata, error)
	Put(context.Context, *model.Metadata) error
}

type Controller struct {
	repo metadataRepository
}

// New creates a new metadata controller.
func New(repo metadataRepository) *Controller {
	return &Controller{repo: repo}
}

// Get returns movie metadata by id.
func (c *Controller) Get(ctx context.Context, id string) (*model.Metadata, error) {
	res, err := c.repo.Get(ctx, id)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return res, nil
}


// Put writes movie metadata to repository.
func (c *Controller) Put(ctx context.Context, m *model.Metadata) error {
	return c.repo.Put(ctx, m)
}