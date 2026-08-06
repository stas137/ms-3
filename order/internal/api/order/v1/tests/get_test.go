package tests

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/stas137/ms-3/order/internal/api/converter"
	"github.com/stas137/ms-3/order/internal/api/order/v1"
	"github.com/stas137/ms-3/order/internal/api/order/v1/mocks"
	errs "github.com/stas137/ms-3/order/internal/errors"
	"github.com/stas137/ms-3/order/internal/model"
	orderv1 "github.com/stas137/ms-3/shared/pkg/openapi/order/v1"
)

func TestGet(t *testing.T) {
	t.Parallel()

	type args struct {
		params orderv1.GetOrderParams
	}

	type expected struct {
		order orderv1.GetOrderRes
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
			{
				PartUUID: uuid.MustParse(gofakeit.UUID()),
				PartType: model.PartTypeHull,
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
		setupMock func(srv *mocks.OrderService)
		expected  expected
	}{
		{
			name: "успешное получение заказа",
			args: args{
				params: orderv1.GetOrderParams{
					OrderUUID: fakeUUID,
				},
			},
			setupMock: func(repo *mocks.OrderService) {
				repo.EXPECT().Get(ctx, fakeUUID).Return(modelOrder, nil)
			},
			expected: expected{
				order: converter.OrderToDTO(modelOrder),
				err:   nil,
			},
		},
		{
			name: "ошибка сервиса при получении заказа (заказ не найден)",
			args: args{
				params: orderv1.GetOrderParams{
					OrderUUID: fakeUUID,
				},
			},
			setupMock: func(repo *mocks.OrderService) {
				repo.EXPECT().Get(ctx, fakeUUID).Return(model.Order{}, errs.ErrOrderNotFound)
			},
			expected: expected{
				order: nil,
				err:   errs.ErrOrderNotFound,
			},
		},
		{
			name: "ошибка сервиса при получении заказа (невалидный uuid)",
			args: args{
				params: orderv1.GetOrderParams{
					OrderUUID: uuid.Nil,
				},
			},
			setupMock: func(repo *mocks.OrderService) {
				// repo.EXPECT().Get(ctx, fakeUUID).Return(model.Order{}, errs.ErrOrderNotFound)
			},
			expected: expected{
				order: nil,
				err:   errs.ErrInvalidUUID,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			orderService := mocks.NewOrderService(t)
			tc.setupMock(orderService)

			api := order.NewApi(orderService)
			res, err := api.GetOrder(ctx, tc.args.params)

			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				assert.Equal(t, nil, res)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected.order, res)
			}
		})
	}
}
