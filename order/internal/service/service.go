package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/cybervasyan/pdididy-project/order/internal/model"
)

type Order interface {
	CreateOrder(ctx context.Context, req *model.Order) (model.Order, error)
	CancelOrder(_ context.Context, orderUUID uuid.UUID) error
	PayOrder(ctx context.Context, req *model.Order) (model.Order, error)
	GetOrderByUuid(_ context.Context, orderUUID uuid.UUID) (model.Order, error)
}
