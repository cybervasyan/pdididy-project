package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/cybervasyan/pdididy-project/inventory/internal/repository/model"
)

type Repository interface {
	Get(ctx context.Context, partUUID uuid.UUID) (model.Part, error)
	List(ctx context.Context, req model.PartsFilter) ([]model.Part, error)
}
