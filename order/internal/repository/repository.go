package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/cybervasyan/pdididy-project/order/internal/repository/model"
)

type Repository interface {
	Create(ctx context.Context, req *model.Order) error
	Cancel(_ context.Context, orderUUID uuid.UUID) error
	Update(ctx context.Context, req *model.Order) (model.Order, error)
	Get(_ context.Context, orderUUID uuid.UUID) (model.Order, error)
}
