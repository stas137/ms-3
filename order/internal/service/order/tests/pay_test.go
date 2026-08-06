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

func TestPay(t *testing.T) {
	t.Parallel()

	type args struct {
		orderUUID uuid.UUID
		method    model.PaymentMethod
	}

	type expected struct {
		// transactionUUID uuid.UUID
		err error
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
		paymentMethod   = model.PaymentMethodCard
		statusPending   = model.OrderStatusPendingPayment
		statusPaid      = model.OrderStatusPaid
		statusCancelled = model.OrderStatusCancelled
		createdAt       = time.Now()
		updatedAt       = time.Now()
		// deletedAt       = time.Now()
		incorrectUUID = ""
	)

	modelOrderGetPending := model.Order{
		UUID:          fakeUUID,
		Items:         items,
		PaymentMethod: &paymentMethod,
		Status:        statusPending,
		CreatedAt:     createdAt,
		UpdatedAt:     &updatedAt,
		DeletedAt:     nil,
	}

	modelOrderGetPaid := model.Order{
		UUID:          fakeUUID,
		Items:         items,
		PaymentMethod: &paymentMethod,
		Status:        statusPaid,
		CreatedAt:     createdAt,
		UpdatedAt:     &updatedAt,
		DeletedAt:     nil,
	}

	modelOrderGetCancelled := model.Order{
		UUID:          fakeUUID,
		Items:         items,
		PaymentMethod: &paymentMethod,
		Status:        statusCancelled,
		CreatedAt:     createdAt,
		UpdatedAt:     &updatedAt,
		DeletedAt:     nil,
	}

	// modelOrderGetDeleted := model.Order{
	// 	UUID:          fakeUUID,
	// 	Items:         items,
	// 	PaymentMethod: &paymentMethod,
	// 	Status:        statusPending,
	// 	CreatedAt:     createdAt,
	// 	UpdatedAt:     &updatedAt,
	// 	DeletedAt:     &deletedAt,
	// }

	tests := []struct {
		name      string
		args      args
		setupMock func(repo *mocks.OrderRepository, paymentClient *mocks.PaymentClient)
		expected  expected
	}{
		{
			name: "успешная оплата заказа",
			args: args{orderUUID: fakeUUID, method: paymentMethod},
			setupMock: func(repo *mocks.OrderRepository, paymentClient *mocks.PaymentClient) {
				repo.On("Get", ctx, fakeUUID).Return(modelOrderGetPending, nil)
				paymentClient.On("PayOrder", ctx, fakeUUID.String(), paymentMethod).Return(fakeUUID.String(), nil)
				repo.On("Update", ctx, mock.MatchedBy(func(modelOrderPay model.Order) bool {
					return (modelOrderPay.UUID == fakeUUID &&
						modelOrderPay.TransactionUUID != nil &&
						*modelOrderPay.PaymentMethod == paymentMethod &&
						modelOrderPay.Status == model.OrderStatusPaid &&
						modelOrderPay.UpdatedAt != nil)
				})).Return(nil)
			},
			expected: expected{
				err: nil,
			},
		},
		{
			name: "ошибка при оплате заказа (заказ не найден)",
			args: args{orderUUID: fakeUUID, method: paymentMethod},
			setupMock: func(repo *mocks.OrderRepository, paymentClient *mocks.PaymentClient) {
				repo.On("Get", ctx, fakeUUID).Return(model.Order{}, errs.ErrOrderNotFound)
			},
			expected: expected{
				err: errs.ErrOrderNotFound,
			},
		},
		{
			name: "ошибка при оплате заказа (заказ уже оплачен)",
			args: args{orderUUID: fakeUUID, method: paymentMethod},
			setupMock: func(repo *mocks.OrderRepository, paymentClient *mocks.PaymentClient) {
				repo.On("Get", ctx, fakeUUID).Return(modelOrderGetPaid, nil)
				// paymentClient.On("PayOrder", ctx, fakeUUID.String(), paymentMethod).Return(uuid.Nil, errs.ErrOrderAlreadyPaid)
			},
			expected: expected{
				err: errs.ErrOrderAlreadyPaid,
			},
		},
		{
			name: "ошибка при оплате заказа (заказ уже отменен)",
			args: args{orderUUID: fakeUUID, method: paymentMethod},
			setupMock: func(repo *mocks.OrderRepository, paymentClient *mocks.PaymentClient) {
				repo.On("Get", ctx, fakeUUID).Return(modelOrderGetCancelled, nil)
				// paymentClient.On("PayOrder", ctx, fakeUUID.String(), paymentMethod).Return(uuid.Nil, errs.)
			},
			expected: expected{
				err: errs.ErrOrderCancelled,
			},
		},
		{
			name: "ошибка при оплате заказа (неверный uuid транзакции заказа)",
			args: args{orderUUID: fakeUUID, method: paymentMethod},
			setupMock: func(repo *mocks.OrderRepository, paymentClient *mocks.PaymentClient) {
				repo.On("Get", ctx, fakeUUID).Return(modelOrderGetPending, nil)
				paymentClient.On("PayOrder", ctx, fakeUUID.String(), paymentMethod).Return(incorrectUUID, nil)
			},
			expected: expected{
				err: errs.ErrInvalidUUID,
			},
		},
		{
			name: "ошибка при оплате заказа (заказ не найден при обновлении заказа)",
			args: args{orderUUID: fakeUUID, method: paymentMethod},
			setupMock: func(repo *mocks.OrderRepository, paymentClient *mocks.PaymentClient) {
				repo.On("Get", ctx, fakeUUID).Return(modelOrderGetPending, nil)
				paymentClient.On("PayOrder", ctx, fakeUUID.String(), paymentMethod).Return(fakeUUID.String(), nil)
				repo.On("Update", ctx, mock.MatchedBy(func(modelOrderPay model.Order) bool {
					return (modelOrderPay.UUID == fakeUUID &&
						modelOrderPay.TransactionUUID != nil &&
						*modelOrderPay.PaymentMethod == paymentMethod &&
						modelOrderPay.Status == model.OrderStatusPaid &&
						modelOrderPay.UpdatedAt != nil)
				})).Return(errs.ErrOrderNotFound)
			},
			expected: expected{
				err: errs.ErrOrderNotFound,
			},
		},
		// {
		// 	name: "ошибка при оплате заказа (заказ не найден при получении Get - заказ удален)",
		// 	args: args{orderUUID: fakeUUID, method: paymentMethod},
		// 	setupMock: func(repo *mocks.OrderRepository, paymentClient *mocks.PaymentClient) {
		// 		repo.On("Get", ctx, fakeUUID).Return(modelOrderGetDeleted, nil)

		// 	},
		// 	expected: expected{
		// 		err: errs.ErrOrderNotFound,
		// 	},
		// },
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			inventoryClient := mocks.NewInventoryClient(t)
			paymentClient := mocks.NewPaymentClient(t)
			orderRepository := mocks.NewOrderRepository(t)

			tc.setupMock(orderRepository, paymentClient)

			svc := order.NewService(orderRepository, inventoryClient, paymentClient)
			res, err := svc.Pay(ctx, tc.args.orderUUID, tc.args.method)

			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				assert.Equal(t, uuid.Nil, res)
			} else {
				require.NoError(t, err)
				assert.NotEqual(t, nil, res)
			}
		})
	}
}
