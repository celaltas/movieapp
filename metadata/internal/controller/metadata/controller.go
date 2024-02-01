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
	}
	return res, nil
}
