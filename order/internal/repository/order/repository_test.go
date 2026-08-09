package order

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/cybervasyan/pdididy-project/order/internal/repository/model"
)

var (
	orderUUID       = uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	userUUID        = uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")
	partUUID        = uuid.MustParse("550e8400-e29b-41d4-a716-446655440002")
	transactionUUID = uuid.MustParse("550e8400-e29b-41d4-a716-446655440003")
)

func testOrder() model.Order {
	method := model.PaymentMethodCARD
	return model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       []uuid.UUID{partUUID},
		TotalPrice:      100,
		TransactionUUID: &transactionUUID,
		PaymentMethod:   &method,
		Status:          model.OrderStatusPAID,
	}
}

func newTestRepository(t *testing.T) (*repository, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	return NewRepository(db), mock
}

func orderRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"order_uuid", "user_uuid", "part_uuids", "total_price",
		"transaction_uuid", "payment_method", "status",
	}).AddRow(
		orderUUID.String(), userUUID.String(), "{"+partUUID.String()+"}", 100.0,
		transactionUUID.String(), string(model.PaymentMethodCARD), string(model.OrderStatusPAID),
	)
}

func TestCreate(t *testing.T) {
	t.Parallel()

	driverErr := errors.New("insert failed")
	tests := []struct {
		name    string
		result  sql.Result
		wantErr error
	}{
		{name: "success", result: sqlmock.NewResult(1, 1)},
		{name: "error", result: sqlmock.NewErrorResult(driverErr), wantErr: driverErr},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo, mock := newTestRepository(t)
			order := testOrder()
			expectation := mock.ExpectExec(regexp.QuoteMeta("INSERT INTO orders("))
			if tt.wantErr != nil {
				expectation.WillReturnError(tt.wantErr)
			} else {
				expectation.WillReturnResult(tt.result)
			}

			err := repo.Create(context.Background(), &order)

			require.ErrorIs(t, err, tt.wantErr)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGet(t *testing.T) {
	t.Parallel()

	driverErr := errors.New("select failed")
	tests := []struct {
		name     string
		rows     *sqlmock.Rows
		queryErr error
		want     model.Order
		wantErr  error
	}{
		{name: "success", rows: orderRows(), want: testOrder()},
		{name: "not found", rows: sqlmock.NewRows([]string{"order_uuid"}), wantErr: model.ErrOrderDoesntExist},
		{name: "error", queryErr: driverErr, wantErr: driverErr},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo, mock := newTestRepository(t)
			expectation := mock.ExpectQuery("SELECT").WithArgs(orderUUID)
			if tt.queryErr != nil {
				expectation.WillReturnError(tt.queryErr)
			} else {
				expectation.WillReturnRows(tt.rows)
			}

			got, err := repo.Get(context.Background(), orderUUID)

			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.want, got)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	driverErr := errors.New("update failed")
	tests := []struct {
		name     string
		rows     *sqlmock.Rows
		queryErr error
		want     model.Order
		wantErr  error
	}{
		{name: "success", rows: orderRows(), want: testOrder()},
		{name: "not found", rows: sqlmock.NewRows([]string{"order_uuid"}), wantErr: model.ErrOrderDoesntExist},
		{name: "error", queryErr: driverErr, wantErr: driverErr},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo, mock := newTestRepository(t)
			order := testOrder()
			expectation := mock.ExpectQuery("UPDATE orders")
			if tt.queryErr != nil {
				expectation.WillReturnError(tt.queryErr)
			} else {
				expectation.WillReturnRows(tt.rows)
			}

			got, err := repo.Update(context.Background(), &order)

			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.want, got)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCancel(t *testing.T) {
	t.Parallel()

	driverErr := errors.New("cancel failed")
	tests := []struct {
		name    string
		result  sql.Result
		execErr error
		wantErr error
	}{
		{name: "success", result: sqlmock.NewResult(0, 1)},
		{name: "not found", result: sqlmock.NewResult(0, 0), wantErr: model.ErrOrderDoesntExist},
		{name: "exec error", execErr: driverErr, wantErr: driverErr},
		{name: "rows affected error", result: sqlmock.NewErrorResult(driverErr), wantErr: driverErr},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo, mock := newTestRepository(t)
			expectation := mock.ExpectExec("UPDATE orders")
			if tt.execErr != nil {
				expectation.WillReturnError(tt.execErr)
			} else {
				expectation.WillReturnResult(tt.result)
			}

			err := repo.Cancel(context.Background(), orderUUID)

			require.ErrorIs(t, err, tt.wantErr)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
