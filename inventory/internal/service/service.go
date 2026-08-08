package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/cybervasyan/pdididy-project/inventory/internal/model"
)

type Part interface {
	GetPart(ctx context.Context, partUUID uuid.UUID) (model.Part, error)
	ListParts(ctx context.Context, req model.PartsFilter) ([]model.Part, error)
}
