package order

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/cybervasyan/pdididy-project/order/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, orderUUID uuid.UUID) (model.Order, error) {
	const query = `
        SELECT
            order_uuid,
            user_uuid,
            part_uuids,
            total_price,
            transaction_uuid,
            payment_method,
            status
        FROM orders
        WHERE order_uuid = $1
    `

	var order model.Order
	var partUUIDs pgtype.FlatArray[pgtype.UUID]
	typeMap := pgtype.NewMap()

	err := r.db.QueryRowContext(ctx, query, orderUUID).Scan(
		&order.OrderUUID,
		&order.UserUUID,
		typeMap.SQLScanner(&partUUIDs),
		&order.TotalPrice,
		&order.TransactionUUID,
		&order.PaymentMethod,
		&order.Status,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return model.Order{}, model.ErrOrderDoesntExist
	}

	if err != nil {
		return model.Order{}, err
	}
	order.PartUuids = make([]uuid.UUID, len(partUUIDs))

	for i, partUUID := range partUUIDs {
		order.PartUuids[i] = uuid.UUID(partUUID.Bytes)
	}

	return order, nil
}
