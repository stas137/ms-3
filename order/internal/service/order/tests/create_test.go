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

	"github.com/stas137/ms-3/order/internal/model"
	"github.com/stas137/ms-3/order/internal/service/input"
	"github.com/stas137/ms-3/order/internal/service/order"
	"github.com/stas137/ms-3/order/internal/service/order/mocks"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type args struct {
		input input.CreateOrderInput
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
		// transactionUUID = uuid.MustParse(gofakeit.UUID())
		// paymentMethod        = model.PaymentMethodCard
		statusPendingPayment = model.OrderStatusPendingPayment
		createdAt            = time.Now()
		// updatedAt            = time.Now()
	)

	modelOrder := model.Order{
		UUID:            fakeUUID,
		Items:           items,
		TransactionUUID: nil,
		// PaymentMethod:   &paymentMethod,
		Status:    statusPendingPayment,
		CreatedAt: createdAt,
		UpdatedAt: nil,
		DeletedAt: nil,
	}

	modelParts := []model.Part{
		{
			UUID:          uuid.New(),
			Name:          "Engine",
			PartType:      model.PartTypeEngine,
			Price:         1000,
			StockQuantity: 10,
		},
		{
			UUID:          uuid.New(),
			Name:          "Hull",
			PartType:      model.PartTypeHull,
			Price:         5000,
			StockQuantity: 5,
		},
	}

	// modelPartsOutOfStock := []model.Part{
	// 	{
	// 		UUID:          uuid.New(),
	// 		Name:          "Engine",
	// 		PartType:      model.PartTypeEngine,
	// 		Price:         1050,
	// 		StockQuantity: 0,
	// 	},
	// }

	tests := []struct {
		name      string
		args      args
		setupMock func(repo *mocks.OrderRepository, inventoryClient *mocks.InventoryClient)
		expected  expected
	}{
		{
			name: "успешное создание заказа",
			args: args{
				input: input.CreateOrderInput{
					HullUUID:   fakeUUID,
					EngineUUID: fakeUUID,
				},
			},
			setupMock: func(repo *mocks.OrderRepository, inventoryClient *mocks.InventoryClient) {
				inventoryClient.On("ListParts", ctx, []string{fakeUUID.String(), fakeUUID.String()}).Return(modelParts, nil)
				inventoryClient.On("ValidateCompatibility", ctx, modelParts[0].UUID, modelParts[1].UUID, (*uuid.UUID)(nil), (*uuid.UUID)(nil)).Return(nil)
				inventoryClient.On("ReserveParts", ctx, []uuid.UUID{modelParts[0].UUID, modelParts[1].UUID}).Return(nil)
				repo.EXPECT().Create(ctx,
					mock.MatchedBy(func(modelOrderCreate model.Order) bool {
						return (modelOrderCreate.UUID != uuid.Nil &&
							// modelOrderCreate.CreatedAt !=  &&
							modelOrderCreate.Status == statusPendingPayment &&
							modelOrderCreate.TransactionUUID == nil)
					}), mock.MatchedBy(func(items []model.OrderItem) bool {
						return (len(items) == 2)
					})).Return(nil)
			},
			expected: expected{
				order: modelOrder,
				err:   nil,
			},
		},
		// {
		// 	name: "ошибка при создании заказа (неверный engine uuid)",
		// 	args: args{input: input.CreateOrderInput{
		// 		HullUUID:   fakeUUID,
		// 		EngineUUID: uuid.Nil,
		// 	}},
		// 	setupMock: func(repo *mocks.OrderRepository, inventoryClient *mocks.InventoryClient) {
		// 	},
		// 	expected: expected{
		// 		order: model.Order{},
		// 		err:   errs.ErrInvalidUUID,
		// 	},
		// },
		// {
		// 	name: "ошибка при создании заказа (неверный hull uuid)",
		// 	args: args{input: input.CreateOrderInput{
		// 		HullUUID:   uuid.Nil,
		// 		EngineUUID: fakeUUID,
		// 	}},
		// 	setupMock: func(repo *mocks.OrderRepository, inventoryClient *mocks.InventoryClient) {
		// 	},
		// 	expected: expected{
		// 		order: model.Order{},
		// 		err:   errs.ErrInvalidUUID,
		// 	},
		// },
		// {
		// 	name: "ошибка при создании заказа (неверный uuid)",
		// 	args: args{input: input.CreateOrderInput{
		// 		HullUUID:   uuid.Nil,
		// 		EngineUUID: uuid.Nil,
		// 	}},
		// 	setupMock: func(repo *mocks.OrderRepository, inventoryClient *mocks.InventoryClient) {
		// 	},
		// 	expected: expected{
		// 		order: model.Order{},
		// 		err:   errs.ErrInvalidUUID,
		// 	},
		// },
		// {
		// 	name: "ошибка при создании заказа (деталь не найдена)",
		// 	args: args{input: input.CreateOrderInput{
		// 		HullUUID:   fakeUUID,
		// 		EngineUUID: fakeUUID,
		// 	}},
		// 	setupMock: func(repo *mocks.OrderRepository, inventoryClient *mocks.InventoryClient) {
		// 		inventoryClient.On("ListParts", ctx, []string{fakeUUID.String(), fakeUUID.String()}).Return([]model.Part{}, errs.ErrPartNotFound)
		// 	},
		// 	expected: expected{
		// 		order: model.Order{},
		// 		err:   errs.ErrPartNotFound,
		// 	},
		// },
		// {
		// 	name: "ошибка при создании заказа (деталь outOfStock)",
		// 	args: args{input: input.CreateOrderInput{
		// 		HullUUID:   fakeUUID,
		// 		EngineUUID: fakeUUID,
		// 	}},
		// 	setupMock: func(repo *mocks.OrderRepository, inventoryClient *mocks.InventoryClient) {
		// 		inventoryClient.On("ListParts", ctx, []string{fakeUUID.String(), fakeUUID.String()}).Return(modelPartsOutOfStock, nil)
		// 	},
		// 	expected: expected{
		// 		order: model.Order{},
		// 		err:   errs.ErrOutOfStock,
		// 	},
		// },
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			orderRepository := mocks.NewOrderRepository(t)
			inventoryClient := mocks.NewInventoryClient(t)

			paymentClient := mocks.NewPaymentClient(t)
			txManager := mocks.NewTxManager(t)

			tc.setupMock(orderRepository, inventoryClient)

			svc := order.NewService(orderRepository, inventoryClient, paymentClient, txManager)
			res, err := svc.Create(ctx, tc.args.input)

			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				assert.Equal(t, model.Order{}, res)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected.order.PaymentMethod, res.PaymentMethod)
				assert.Equal(t, tc.expected.order.Status, res.Status)
				assert.Equal(t, tc.expected.order.TransactionUUID, res.TransactionUUID)
			}
		})
	}
}
