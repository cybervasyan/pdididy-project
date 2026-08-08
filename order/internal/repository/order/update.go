package order

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/cybervasyan/pdididy-project/order/internal/repository/model"
)

func (r *repository) Update(ctx context.Context, req *model.Order) (model.Order, error) {
	const query = `
        UPDATE orders
        SET
            user_uuid = $2,
            part_uuids = $3,
            total_price = $4,
            transaction_uuid = $5,
            payment_method = $6,
            status = $7
        WHERE order_uuid = $1
        RETURNING
            order_uuid,
            user_uuid,
            part_uuids,
            total_price,
            transaction_uuid,
            payment_method,
            status
    `

	var updated model.Order
	var partUUIDs pgtype.FlatArray[pgtype.UUID]
	typeMap := pgtype.NewMap()

	err := r.db.QueryRowContext(
		ctx,
		query,
		req.OrderUUID,
		req.UserUUID,
		req.PartUuids,
		req.TotalPrice,
		req.TransactionUUID,
		req.PaymentMethod,
		req.Status,
	).Scan(
		&updated.OrderUUID,
		&updated.UserUUID,
		typeMap.SQLScanner(&partUUIDs),
		&updated.TotalPrice,
		&updated.TransactionUUID,
		&updated.PaymentMethod,
		&updated.Status,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return model.Order{}, model.ErrOrderDoesntExist
	}

	if err != nil {
		return model.Order{}, err
	}

	updated.PartUuids = make([]uuid.UUID, len(partUUIDs))

	for i, partUUID := range partUUIDs {
		updated.PartUuids[i] = uuid.UUID(partUUID.Bytes)
	}

	return updated, nil
}
