package tests

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	errs "github.com/stas137/ms-3/order/internal/errors"
	"github.com/stas137/ms-3/order/internal/model"
	"github.com/stas137/ms-3/order/internal/service/order"
	"github.com/stas137/ms-3/order/internal/service/order/mocks"
)

func TestCancel(t *testing.T) {
	t.Parallel()

	type args struct {
		orderUUID uuid.UUID
	}

	type expected struct {
		err error
	}

	var (
		ctx      = context.Background()
		fakeUUID = uuid.New()
		items    = []model.OrderItem{
			{
				PartUUID: uuid.MustParse(gofakeit.UUID()),
				PartType: model.PartTypeEngine,
				Price:    int64(gofakeit.Price(100, 100000)),
			},
		}
		paymentMethod   = model.PaymentMethodCard
		status          = model.OrderStatusPendingPayment
		statusCancelled = model.OrderStatusCancelled
		createdAt       = time.Now()
	)

	modelOrder := model.Order{
		UUID:            fakeUUID,
		Items:           items,
		TransactionUUID: nil,
		PaymentMethod:   &paymentMethod,
		Status:          status,
		CreatedAt:       createdAt,
		UpdatedAt:       nil,
		DeletedAt:       nil,
	}

	// modelOrderCancelled := model.Order{
	// 	UUID:            fakeUUID,
	// 	Items:           items,
	// 	TransactionUUID: nil,
	// 	PaymentMethod:   &paymentMethod,
	// 	Status:          statusCancelled,
	// 	CreatedAt:       createdAt,
	// 	UpdatedAt:       &time.Time{},
	// 	DeletedAt:       nil,
	// }

	tests := []struct {
		name      string
		args      args
		setupMock func(repo *mocks.OrderRepository)
		expected  expected
	}{
		{
			name: "успешная отмена заказа",
			args: args{orderUUID: fakeUUID},
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().Get(ctx, fakeUUID).Return(modelOrder, nil)
				repo.EXPECT().Update(ctx, mock.MatchedBy(func(modelOrder model.Order) bool {
					return (modelOrder.UUID == fakeUUID &&
						modelOrder.TransactionUUID == nil &&
						*modelOrder.PaymentMethod == paymentMethod &&
						modelOrder.Status == statusCancelled &&
						modelOrder.UpdatedAt != nil)
				})).Return(nil)
			},
			expected: expected{
				err: nil,
			},
		},
		{
			name: "ошибка при отмене заказа (заказ уже отменен)",
			args: args{orderUUID: fakeUUID},
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().Get(ctx, fakeUUID).Return(model.Order{}, errs.ErrOrderCancelled)
			},
			expected: expected{
				err: errs.ErrOrderCancelled,
			},
		},
		{
			name: "ошибка при отмене заказа (заказ уже оплачен)",
			args: args{orderUUID: fakeUUID},
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().Get(ctx, fakeUUID).Return(model.Order{}, errs.ErrOrderAlreadyPaid)
			},
			expected: expected{
				err: errs.ErrOrderAlreadyPaid,
			},
		},
		{
			name: "ошибка при отмене заказа (заказ не найден)",
			args: args{orderUUID: fakeUUID},
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().Get(ctx, fakeUUID).Return(model.Order{}, errs.ErrOrderNotFound)
			},
			expected: expected{
				err: errs.ErrOrderNotFound,
			},
		},
		{
			name: "ошибка при отмене заказа (заказ не найден при обновлении)",
			args: args{orderUUID: fakeUUID},
			setupMock: func(repo *mocks.OrderRepository) {
				repo.EXPECT().Get(ctx, fakeUUID).Return(modelOrder, nil)
				repo.EXPECT().Update(ctx, mock.MatchedBy(func(modelOrder model.Order) bool {
					return (modelOrder.UUID == fakeUUID &&
						modelOrder.TransactionUUID == nil &&
						*modelOrder.PaymentMethod == paymentMethod &&
						modelOrder.Status == model.OrderStatusCancelled &&
						modelOrder.UpdatedAt != nil)
				})).Return(errs.ErrOrderNotFound)
			},
			expected: expected{
				err: errs.ErrOrderNotFound,
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
			err := svc.Cancel(ctx, tc.args.orderUUID)

			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
