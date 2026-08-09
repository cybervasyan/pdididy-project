package order

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/cybervasyan/pdididy-project/order/internal/model"
	repoModel "github.com/cybervasyan/pdididy-project/order/internal/repository/model"
	inventoryv1 "github.com/cybervasyan/pdididy-project/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/cybervasyan/pdididy-project/shared/pkg/proto/payment/v1"
)

type repositoryStub struct {
	create func(context.Context, *repoModel.Order) error
	get    func(context.Context, uuid.UUID) (repoModel.Order, error)
	update func(context.Context, *repoModel.Order) (repoModel.Order, error)
	cancel func(context.Context, uuid.UUID) error
}

func (s repositoryStub) Create(ctx context.Context, order *repoModel.Order) error {
	return s.create(ctx, order)
}

func (s repositoryStub) Get(ctx context.Context, id uuid.UUID) (repoModel.Order, error) {
	return s.get(ctx, id)
}

func (s repositoryStub) Update(ctx context.Context, order *repoModel.Order) (repoModel.Order, error) {
	return s.update(ctx, order)
}

func (s repositoryStub) Cancel(ctx context.Context, id uuid.UUID) error {
	return s.cancel(ctx, id)
}

type inventoryStub struct {
	list func(context.Context, *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error)
}

func (s inventoryStub) ListParts(ctx context.Context, request *inventoryv1.ListPartsRequest, _ ...grpc.CallOption) (*inventoryv1.ListPartsResponse, error) {
	return s.list(ctx, request)
}

type paymentStub struct {
	pay func(context.Context, *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error)
}

func (s paymentStub) PayOrder(ctx context.Context, request *paymentv1.PayOrderRequest, _ ...grpc.CallOption) (*paymentv1.PayOrderResponse, error) {
	return s.pay(ctx, request)
}

func unusedRepository() repositoryStub {
	return repositoryStub{
		create: func(context.Context, *repoModel.Order) error { panic("unexpected Create") },
		get:    func(context.Context, uuid.UUID) (repoModel.Order, error) { panic("unexpected Get") },
		update: func(context.Context, *repoModel.Order) (repoModel.Order, error) { panic("unexpected Update") },
		cancel: func(context.Context, uuid.UUID) error { panic("unexpected Cancel") },
	}
}

func unusedInventory() inventoryStub {
	return inventoryStub{list: func(context.Context, *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
		panic("unexpected ListParts")
	}}
}

func unusedPayment() paymentStub {
	return paymentStub{pay: func(context.Context, *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
		panic("unexpected PayOrder")
	}}
}

func TestCreateOrder(t *testing.T) {
	t.Parallel()

	partID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440010")
	userID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440011")
	driverErr := errors.New("dependency failed")
	request := model.Order{UserUUID: userID, PartUuids: []uuid.UUID{partID}}

	tests := []struct {
		name    string
		repo    repositoryStub
		inv     inventoryStub
		wantErr error
	}{
		{
			name: "success",
			repo: func() repositoryStub {
				repo := unusedRepository()
				repo.create = func(_ context.Context, order *repoModel.Order) error {
					require.NotEqual(t, uuid.Nil, order.OrderUUID)
					require.Equal(t, 42.5, order.TotalPrice)
					require.Equal(t, repoModel.OrderStatusPENDINGPAYMENT, order.Status)
					return nil
				}
				return repo
			}(),
			inv: inventoryStub{list: func(_ context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
				require.Equal(t, []string{partID.String()}, req.Filter.Uuids)
				return &inventoryv1.ListPartsResponse{Parts: []*inventoryv1.Part{{Uuid: partID.String(), Price: 42.5}}}, nil
			}},
		},
		{
			name: "inventory error",
			repo: unusedRepository(),
			inv: inventoryStub{list: func(context.Context, *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
				return nil, driverErr
			}},
			wantErr: driverErr,
		},
		{
			name: "missing part",
			repo: unusedRepository(),
			inv: inventoryStub{list: func(context.Context, *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
				return &inventoryv1.ListPartsResponse{}, nil
			}},
			wantErr: model.ErrPartNotFound,
		},
		{
			name: "repository error",
			repo: func() repositoryStub {
				repo := unusedRepository()
				repo.create = func(context.Context, *repoModel.Order) error { return driverErr }
				return repo
			}(),
			inv: inventoryStub{list: func(context.Context, *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
				return &inventoryv1.ListPartsResponse{Parts: []*inventoryv1.Part{{Uuid: partID.String(), Price: 1}}}, nil
			}},
			wantErr: driverErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			service := NewOrderService(tt.repo, tt.inv, unusedPayment())

			got, err := service.CreateOrder(context.Background(), &request)

			require.ErrorIs(t, err, tt.wantErr)
			if tt.wantErr == nil {
				require.Equal(t, userID, got.UserUUID)
				require.Equal(t, 42.5, got.TotalPrice)
			}
		})
	}
}

func TestGetOrderByUUID(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	driverErr := errors.New("get failed")
	method := repoModel.PaymentMethodCARD
	transactionID := uuid.New()
	want := repoModel.Order{OrderUUID: id, Status: repoModel.OrderStatusPAID, PaymentMethod: &method, TransactionUUID: &transactionID}

	tests := []struct {
		name    string
		result  repoModel.Order
		err     error
		wantErr error
	}{
		{name: "success", result: want},
		{name: "not found", err: repoModel.ErrOrderDoesntExist, wantErr: model.ErrOrderDoesntExist},
		{name: "repository error", err: driverErr, wantErr: driverErr},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := unusedRepository()
			repo.get = func(context.Context, uuid.UUID) (repoModel.Order, error) { return tt.result, tt.err }
			service := NewOrderService(repo, unusedInventory(), unusedPayment())

			got, err := service.GetOrderByUuid(context.Background(), id)

			require.ErrorIs(t, err, tt.wantErr)
			if tt.wantErr == nil {
				require.Equal(t, model.OrderStatusPAID, got.Status)
				require.Equal(t, model.PaymentMethodCARD, *got.PaymentMethod)
				require.Equal(t, transactionID, *got.TransactionUUID)
			}
		})
	}
}

func TestCancelOrder(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	driverErr := errors.New("cancel failed")

	tests := []struct {
		name      string
		getResult repoModel.Order
		getErr    error
		cancelErr error
		wantErr   error
	}{
		{name: "success", getResult: repoModel.Order{Status: repoModel.OrderStatusPENDINGPAYMENT}},
		{name: "not found", getErr: repoModel.ErrOrderDoesntExist, wantErr: model.ErrOrderDoesntExist},
		{name: "get error", getErr: driverErr, wantErr: driverErr},
		{name: "wrong status", getResult: repoModel.Order{Status: repoModel.OrderStatusPAID}, wantErr: model.ErrOrderNotInPending},
		{name: "cancel error", getResult: repoModel.Order{Status: repoModel.OrderStatusPENDINGPAYMENT}, cancelErr: driverErr, wantErr: driverErr},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := unusedRepository()
			repo.get = func(context.Context, uuid.UUID) (repoModel.Order, error) { return tt.getResult, tt.getErr }
			repo.cancel = func(context.Context, uuid.UUID) error { return tt.cancelErr }
			service := NewOrderService(repo, unusedInventory(), unusedPayment())

			err := service.CancelOrder(context.Background(), id)

			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestPayOrder(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	userID := uuid.New()
	transactionID := uuid.New()
	method := model.PaymentMethodCARD
	driverErr := errors.New("payment failed")
	request := model.Order{OrderUUID: id, PaymentMethod: &method}
	pending := repoModel.Order{OrderUUID: id, UserUUID: userID, Status: repoModel.OrderStatusPENDINGPAYMENT}

	tests := []struct {
		name        string
		getResult   repoModel.Order
		getErr      error
		paymentErr  error
		transaction string
		updateErr   error
		wantErr     error
		wantAnyErr  bool
	}{
		{name: "success", getResult: pending, transaction: transactionID.String()},
		{name: "not found", getErr: repoModel.ErrOrderDoesntExist, wantErr: model.ErrOrderDoesntExist},
		{name: "get error", getErr: driverErr, wantErr: driverErr},
		{name: "wrong status", getResult: repoModel.Order{Status: repoModel.OrderStatusPAID}, wantErr: model.ErrOrderNotInPending},
		{name: "payment error", getResult: pending, paymentErr: driverErr, wantErr: driverErr},
		{name: "invalid transaction", getResult: pending, transaction: "invalid", wantAnyErr: true},
		{name: "update error", getResult: pending, transaction: transactionID.String(), updateErr: driverErr, wantErr: driverErr},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := unusedRepository()
			repo.get = func(context.Context, uuid.UUID) (repoModel.Order, error) { return tt.getResult, tt.getErr }
			repo.update = func(_ context.Context, order *repoModel.Order) (repoModel.Order, error) {
				return *order, tt.updateErr
			}
			payment := paymentStub{pay: func(_ context.Context, req *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
				require.Equal(t, id.String(), req.OrderUuid)
				return &paymentv1.PayOrderResponse{TransactionUuid: tt.transaction}, tt.paymentErr
			}}
			service := NewOrderService(repo, unusedInventory(), payment)

			got, err := service.PayOrder(context.Background(), &request)

			if tt.wantAnyErr {
				require.Error(t, err)
			} else {
				require.ErrorIs(t, err, tt.wantErr)
			}
			if tt.wantErr == nil && !tt.wantAnyErr {
				require.Equal(t, model.OrderStatusPAID, got.Status)
				require.Equal(t, transactionID, *got.TransactionUUID)
			}
		})
	}
}
