package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/stas137/ms-3/order/internal/model"
)

// OrderRepository определяет контракт для работы с хранилищем заказов
// Интерфейс определен по месту использования - в пакете service
type OrderRepository interface {
	Create(ctx context.Context, order model.Order) error
	Get(ctx context.Context, orderUUID uuid.UUID) (model.Order, error)
	Update(ctx context.Context, order model.Order) error
}

type InventoryClient interface {
	ListParts(ctx context.Context, uuids []string) ([]model.Part, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID string, method model.PaymentMethod) (string, error)
}
