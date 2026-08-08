package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/cybervasyan/pdididy-project/order/internal/repository/model"
)

func (r *repository) Cancel(ctx context.Context, orderUUID uuid.UUID) error {
	const query = `
		UPDATE orders
		SET status = $2
		WHERE order_uuid = $1
	`
	result, err := r.db.ExecContext(
		ctx,
		query,
		orderUUID,
		model.OrderStatusCANCELLED,
	)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return model.ErrOrderDoesntExist
	}

	return nil
}
