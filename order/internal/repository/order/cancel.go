package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/cybervasyan/pdididy-project/order/internal/repository/model"
)

func (r *repository) Cancel(_ context.Context, orderUUID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, ok := r.orders[orderUUID]
	if !ok {
		return model.ErrOrderDoesntExist
	}

	order.Status = model.OrderStatusCANCELLED

	return nil
}
