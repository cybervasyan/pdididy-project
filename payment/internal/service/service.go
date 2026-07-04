package service

import (
	"context"

	"github.com/cybervasyan/pdididy-project/payment/internal/model"
	"github.com/google/uuid"
)

type Payment interface {
	PayOrder(ctx context.Context, orderUUID, userUUID uuid.UUID, paymentMethod model.PaymentMethod) (uuid.UUID, error)
}
