package handler

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	orderv1 "github.com/stas137/ms-3/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/stas137/ms-3/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/stas137/ms-3/shared/pkg/proto/payment/v1"
)

// OrderStatus — статус заказа
type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

// PaymentMethod — способ оплаты заказа
type PaymentMethod string

const (
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

// Order представляет заказ на постройку космического корабля
type Order struct {
	OrderUUID       uuid.UUID
	HullUUID        uuid.UUID
	EngineUUID      uuid.UUID
	ShieldUUID      *uuid.UUID // опциональный
	WeaponUUID      *uuid.UUID // опциональный
	TotalPrice      int64      // в копейках
	TransactionUUID *uuid.UUID
	PaymentMethod   *PaymentMethod
	Status          OrderStatus
	CreatedAt       time.Time
}

// orderStore — хранилище заказов (in-memory) потокобезовасное
type orderStore struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]Order
}

// NewOrderStore создаёт новое пустое хранилище заказов
func NewOrderStore() *orderStore {
	return &orderStore{
		orders: make(map[uuid.UUID]Order),
	}
}

// handler реализует интерфейс orderv1.Handler, сгенерированный ogen
type handler struct {
	orderv1.UnimplementedHandler
	inventoryClient inventoryv1.InventoryServiceClient
	paymentClient   paymentv1.PaymentServiceClient
	store           *orderStore
}

// NewHandler создаёт новый обработчик заказов
func NewHandler(
	inventoryClient inventoryv1.InventoryServiceClient,
	paymentClient paymentv1.PaymentServiceClient,
	store *orderStore,
) *handler {
	return &handler{
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		store:           store,
	}
}

// SetupServer создаёт OpenAPI сервер на основе обработчика
func SetupServer(h *handler) (*orderv1.Server, error) {
	return orderv1.NewServer(h)
}

func (h *orderStore) GetOrder(orderUUID uuid.UUID) (Order, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	order, ok := h.orders[orderUUID]

	return order, ok
}

func (h *orderStore) CreateOrder(createOrder Order) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.orders[createOrder.OrderUUID] = createOrder
}

func (h *orderStore) PayOrder(payOrder Order) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.orders[payOrder.OrderUUID] = payOrder
}

func (h *orderStore) CancelOrder(cancelOrder Order) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.orders[cancelOrder.OrderUUID] = cancelOrder
}

// GetOrder реализует операцию getOrder (пример реализации)
// GET /api/v1/orders/{order_uuid}.
func (h *handler) GetOrder(_ context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	// 1. Найти заказ в store (с блокировкой для thread-safety)
	h.store.mu.RLock()
	order, ok := h.store.orders[params.OrderUUID]
	h.store.mu.RUnlock()

	// 2. Если не найден — вернуть 404
	if !ok {
		return &orderv1.GetOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	}

	// 3. Преобразовать в DTO и вернуть
	var shieldUUID orderv1.OptNilUUID
	if order.ShieldUUID != nil {
		shieldUUID = orderv1.NewOptNilUUID(*order.ShieldUUID)
	}

	var weaponUUID orderv1.OptNilUUID
	if order.WeaponUUID != nil {
		weaponUUID = orderv1.NewOptNilUUID(*order.WeaponUUID)
	}

	var transactionUUID orderv1.OptNilUUID
	if order.TransactionUUID != nil {
		transactionUUID = orderv1.NewOptNilUUID(*order.TransactionUUID)
	}

	var paymentMethod orderv1.OptNilPaymentMethod
	if order.PaymentMethod != nil {
		paymentMethod = orderv1.NewOptNilPaymentMethod(orderv1.PaymentMethod(*order.PaymentMethod))
	}

	return &orderv1.OrderDto{
		OrderUUID:       order.OrderUUID,
		HullUUID:        order.HullUUID,
		EngineUUID:      order.EngineUUID,
		ShieldUUID:      shieldUUID,
		WeaponUUID:      weaponUUID,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          orderv1.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
	}, nil
}

// Реализовать остальные методы интерфейса orderv1.Handler:
//
// CreateOrder реализует операцию createOrder
// POST /api/v1/orders
func (h *handler) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	// 1. Валидация: hull_uuid и engine_uuid обязательны
	// 2. Получить детали через InventoryService.ListParts
	// 3. Проверить stock_quantity > 0
	// 4. Вычислить total_price
	// 5. Сгенерировать order_uuid (UUID v4)
	// 6. Создать заказ со статусом PENDING_PAYMENT
	// 7. Сохранить в store
	// 8. Вернуть order_uuid и total_price

	engineUUID := req.GetEngineUUID()
	if engineUUID.String() == "" {
		return &orderv1.CreateOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "uuid двигателя обязательный",
		}, nil
	}

	if _, err := uuid.Parse(engineUUID.String()); err != nil {
		return &orderv1.CreateOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "неверный формат uuid двигателя",
		}, nil
	}

	hullUUID := req.GetHullUUID()
	if hullUUID.String() == "" {
		return &orderv1.CreateOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "uuid корпуса обязательный",
		}, nil
	}

	if _, err := uuid.Parse(hullUUID.String()); err != nil {
		return &orderv1.CreateOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "неверный формат uuid корпуса",
		}, nil
	}

	uuids := []string{
		engineUUID.String(),
		hullUUID.String(),
	}

	var weaponUUID uuid.UUID
	if req.GetWeaponUUID().Value != uuid.Nil {
		weaponUUID = req.GetWeaponUUID().Value
		uuids = append(uuids, weaponUUID.String())
	}

	var shieldUUID uuid.UUID
	if req.GetShieldUUID().Value != uuid.Nil {
		shieldUUID = req.GetShieldUUID().Value
		uuids = append(uuids, shieldUUID.String())
	}

	listParts, err := h.inventoryClient.ListParts(
		ctx,
		&inventoryv1.ListPartsRequest{
			Uuids: uuids,
		},
	)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return &orderv1.CreateOrderInternalServerError{
				Code:    http.StatusInternalServerError,
				Message: "внутренняя ошибка сервера",
			}, nil
		}

		switch st.Code() {
		case codes.NotFound:
			return &orderv1.CreateOrderNotFound{
				Code:    http.StatusNotFound,
				Message: " деталь с таким uuid не найдена",
			}, nil
		default:
			return &orderv1.CreateOrderInternalServerError{
				Code:    http.StatusInternalServerError,
				Message: "внутренняя ошибка сервера - не удалось получить детали заказа",
			}, nil
		}

	}

	parts := listParts.GetParts()
	var totalPrice int64

	for _, part := range parts {
		if part.GetStockQuantity() == 0 {
			return &orderv1.CreateOrderConflict{
				Code:    http.StatusConflict,
				Message: "не найдены указанные детали",
			}, nil
		}

		totalPrice += part.GetPrice()
	}

	orderUUID := uuid.New()

	order := Order{
		OrderUUID:  orderUUID,
		HullUUID:   hullUUID,
		EngineUUID: engineUUID,
		WeaponUUID: &weaponUUID,
		ShieldUUID: &shieldUUID,
		TotalPrice: totalPrice,
		Status:     OrderStatusPendingPayment,
		CreatedAt:  time.Now(),
	}

	h.store.CreateOrder(order)

	return &orderv1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: totalPrice,
	}, nil
}

// PayOrder реализует операцию payOrder
// POST /api/v1/orders/{order_uuid}/pay
func (h *handler) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	// 1. Найти заказ в store
	// 2. Проверить статус == PENDING_PAYMENT
	// 3. Вызвать h.paymentClient.PayOrder для обработки платежа
	// 4. Обновить статус на PAID и сохранить transaction_uuid
	// 5. Вернуть transaction_uuid

	orderUUID := params.OrderUUID

	order, ok := h.store.GetOrder(orderUUID)
	if !ok {
		return &orderv1.PayOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	}

	orderv1PaymentMethod := req.GetPaymentMethod()
	paymentMethod := "PAYMENT_METHOD_" + string(orderv1PaymentMethod)
	paymentv1PaymentMethod := paymentv1.PaymentMethod(paymentv1.PaymentMethod_value[paymentMethod])

	if order.Status == OrderStatusPendingPayment {
		payResult, err := h.paymentClient.PayOrder(ctx, &paymentv1.PayOrderRequest{
			OrderUuid:     orderUUID.String(),
			PaymentMethod: paymentv1PaymentMethod,
		})
		if err != nil {
			return &orderv1.PayOrderInternalServerError{
				Code:    http.StatusInternalServerError,
				Message: "внутренняя ошибка сервера",
			}, nil
		}

		transactionUUID, err := uuid.Parse(payResult.TransactionUuid)
		if err != nil {
			return &orderv1.PayOrderInternalServerError{
				Code:    http.StatusInternalServerError,
				Message: "внутренняя ошибка сервера",
			}, nil
		}

		orderPaymentMethod := PaymentMethod(req.GetPaymentMethod())

		order.Status = OrderStatusPaid
		order.TransactionUUID = &transactionUUID
		order.PaymentMethod = &orderPaymentMethod // переписать в одну строку

		h.store.PayOrder(order)

		return &orderv1.PayOrderResponse{
			TransactionUUID: transactionUUID,
		}, nil
	}

	if order.Status == OrderStatusPaid {
		return &orderv1.PayOrderConflict{
			Code:    http.StatusConflict,
			Message: "заказ уже оплачен",
		}, nil
	}

	if order.Status == OrderStatusCancelled {
		return &orderv1.PayOrderConflict{
			Code:    http.StatusConflict,
			Message: "заказ отменен",
		}, nil
	}

	return nil, errors.New("неизвестная ошибка")
}

// CancelOrder реализует операцию cancelOrder
// POST /api/v1/orders/{order_uuid}/cancel
func (h *handler) CancelOrder(ctx context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
	// 1. Найти заказ в store
	// 2. Проверить статус == PENDING_PAYMENT
	// 3. Обновить статус на CANCELLED
	// 4. Вернуть success

	orderUUID := params.OrderUUID

	cancelOrder, ok := h.store.orders[orderUUID]
	if !ok {
		return &orderv1.CancelOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	}

	if cancelOrder.Status == OrderStatusPendingPayment {
		cancelOrder.Status = OrderStatusCancelled

		h.store.CancelOrder(cancelOrder)

		return &orderv1.CancelOrderResponse{}, nil
	}

	return &orderv1.CancelOrderConflict{
		Code:    http.StatusConflict,
		Message: "order conflict",
	}, nil
}
