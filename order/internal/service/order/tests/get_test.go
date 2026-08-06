package tests

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/stas137/ms-3/order/internal/errors"
	"github.com/stas137/ms-3/order/internal/model"
	"github.com/stas137/ms-3/order/internal/service/order"
	"github.com/stas137/ms-3/order/internal/service/order/mocks"
)

func TestGet(t *testing.T) {
	t.Parallel()

	type args struct {
		orderUUID uuid.UUID
	}

	type expected struct {
		order model.Order
		err   error
	}

	var (
		ctx      = context.Background()
		fakeUUID = uuid.MustParse(gofakeit.UUID())
		items    = []model.OrderItem{
			{
				PartUUID: uuid.MustParse(gofakeit.UUID()),
				PartType: model.PartTypeEngine,
				Price:    int64(gofakeit.Price(100, 100000)),
			},
		}
		transactionUUID = uuid.MustParse(gofakeit.UUID())
		paymentMethod   = model.PaymentMethodCard
		status          = model.OrderStatusPaid
		createdAt       = time.Now()
		updatedAt       = time.Now()
	)

	modelOrder := model.Order{
		UUID:            fakeUUID,
		Items:           items,
		TransactionUUID: &transactionUUID,
		PaymentMethod:   &paymentMethod,
		Status:          status,
		CreatedAt:       createdAt,
		UpdatedAt:       &updatedAt,
		DeletedAt:       nil,
	}

	tests := []struct {
		name      string
		args      args
		setupMock func(repo *mocks.OrderRepository)
		expected  expected
	}{
		{
			name: "успешное получение заказа",
			args: args{orderUUID: fakeUUID},
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().Get(ctx, fakeUUID).Return(modelOrder, nil)
			},
			expected: expected{
				order: modelOrder,
				err:   nil,
			},
		}, {
			name: "ошибка репозитория при получении заказа",
			args: args{orderUUID: fakeUUID},
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().Get(ctx, fakeUUID).Return(model.Order{}, errs.ErrOrderNotFound)
			},
			expected: expected{
				order: model.Order{},
				err:   errs.ErrOrderNotFound,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			orderRepository := mocks.NewOrderRepository(t)
			tc.setupMock(orderRepository)

			inventoryClient := mocks.NewInventoryClient(t)
			paymentClient := mocks.NewPaymentClient(t)

			svc := order.NewService(orderRepository, inventoryClient, paymentClient)
			res, err := svc.Get(ctx, tc.args.orderUUID)

			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				assert.Equal(t, model.Order{}, res)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected.order, res)
			}
		})
	}
}
