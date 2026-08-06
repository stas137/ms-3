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
	"github.com/stas137/ms-3/order/internal/model"
	"github.com/stas137/ms-3/order/internal/service/input"
	orderv1 "github.com/stas137/ms-3/shared/pkg/openapi/order/v1"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type args struct {
		req *orderv1.CreateOrderRequest
	}

	type expected struct {
		createRes orderv1.CreateOrderRes
		err       error
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
			name: "успешное создание заказа",
			args: args{
				req: &orderv1.CreateOrderRequest{
					HullUUID:   fakeUUID,
					EngineUUID: fakeUUID,
				},
			},
			setupMock: func(repo *mocks.OrderService) {
				repo.EXPECT().Create(ctx, input.CreateOrderInput{HullUUID: fakeUUID, EngineUUID: fakeUUID}).Return(model.Order{
					UUID: fakeUUID,
					Items: []model.OrderItem{
						{
							PartUUID: fakeUUID,
							PartType: model.PartTypeHull,
							Price:    2500,
						},
						{
							PartUUID: fakeUUID,
							PartType: model.PartTypeEngine,
							Price:    2500,
						},
					},
				}, nil)
			},
			expected: expected{
				createRes: &orderv1.CreateOrderResponse{
					OrderUUID:  fakeUUID,
					TotalPrice: 5000,
				},
				err: nil,
			},
		},
		{
			name: "ошибка сервиса при создании заказа (невалидный uuid)",
			args: args{
				req: &orderv1.CreateOrderRequest{
					HullUUID:   fakeUUID,
					EngineUUID: uuid.Nil,
				},
			},
			setupMock: func(repo *mocks.OrderService) {
				repo.EXPECT().Create(ctx, input.CreateOrderInput{
					HullUUID:   fakeUUID,
					EngineUUID: uuid.Nil,
				}).Return(model.Order{}, errs.ErrInvalidUUID)
			},
			expected: expected{
				createRes: &orderv1.CreateOrderResponse{},
				err:       errs.ErrInvalidUUID,
			},
		},
		// {
		// 	name: "ошибка сервиса при получении заказа (невалидный uuid)",
		// 	args: args{
		// 		params: orderv1.CancelOrderParams{
		// 			OrderUUID: uuid.Nil,
		// 		},
		// 	},
		// 	setupMock: func(repo *mocks.OrderService) {
		// 		repo.EXPECT().Cancel(ctx, uuid.Nil).Return(errs.ErrInvalidUUID)
		// 	},
		// 	expected: expected{
		// 		order: &orderv1.CancelOrderResponse{},
		// 		err:   errs.ErrInvalidUUID,
		// 	},
		// },
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			orderService := mocks.NewOrderService(t)
			tc.setupMock(orderService)

			api := order.NewApi(orderService)
			res, err := api.CreateOrder(ctx, tc.args.req)

			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				assert.Equal(t, &orderv1.CreateOrderResponse{}, res)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected.createRes, res)
			}
		})
	}
}
