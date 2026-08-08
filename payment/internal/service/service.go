package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/cybervasyan/pdididy-project/payment/internal/model"
)

type Payment interface {
	PayOrder(ctx context.Context, orderUUID, userUUID uuid.UUID, paymentMethod model.PaymentMethod) (uuid.UUID, error)
}
