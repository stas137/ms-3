package tests

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/stas137/ms-3/order/internal/api/order/v1"
	"github.com/stas137/ms-3/order/internal/api/order/v1/mocks"
	errs "github.com/stas137/ms-3/order/internal/errors"
	orderv1 "github.com/stas137/ms-3/shared/pkg/openapi/order/v1"
)

func TestCancel(t *testing.T) {
	t.Parallel()

	type args struct {
		params orderv1.CancelOrderParams
	}

	type expected struct {
		order orderv1.CancelOrderRes
		err   error
	}

	var (
		ctx      = context.Background()
		fakeUUID = uuid.MustParse(gofakeit.UUID())
	)

	tests := []struct {
		name      string
		args      args
		setupMock func(srv *mocks.OrderService)
		expected  expected
	}{
		{
			name: "успешная отмена заказа",
			args: args{
				params: orderv1.CancelOrderParams{
					OrderUUID: fakeUUID,
				},
			},
			setupMock: func(repo *mocks.OrderService) {
				repo.EXPECT().Cancel(ctx, fakeUUID).Return(nil)
			},
			expected: expected{
				order: &orderv1.CancelOrderResponse{},
				err:   nil,
			},
		},
		{
			name: "ошибка сервиса при получении заказа (заказ не найден)",
			args: args{
				params: orderv1.CancelOrderParams{
					OrderUUID: fakeUUID,
				},
			},
			setupMock: func(repo *mocks.OrderService) {
				repo.EXPECT().Cancel(ctx, fakeUUID).Return(errs.ErrOrderNotFound)
			},
			expected: expected{
				order: &orderv1.CancelOrderResponse{},
				err:   errs.ErrOrderNotFound,
			},
		},
		{
			name: "ошибка сервиса при получении заказа (невалидный uuid)",
			args: args{
				params: orderv1.CancelOrderParams{
					OrderUUID: uuid.Nil,
				},
			},
			setupMock: func(repo *mocks.OrderService) {
				repo.EXPECT().Cancel(ctx, uuid.Nil).Return(errs.ErrInvalidUUID)
			},
			expected: expected{
				order: &orderv1.CancelOrderResponse{},
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
			res, err := api.CancelOrder(ctx, tc.args.params)

			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				assert.Equal(t, &orderv1.CancelOrderResponse{}, res)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected.order, res)
			}
		})
	}
}
