package order

import (
	"context"

	"github.com/cybervasyan/pdididy-project/order/internal/repository/model"
)

func (r *repository) Create(ctx context.Context, req *model.Order) error {
	const query = `
		INSERT INTO orders(
		    order_uuid,
            user_uuid,
            part_uuids,
            total_price,
            transaction_uuid,
            payment_method,
            status
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		req.OrderUUID,
		req.UserUUID,
		uuidArray(req.PartUuids),
		req.TotalPrice,
		req.TransactionUUID,
		req.PaymentMethod,
		req.Status,
	)

	return err
}
