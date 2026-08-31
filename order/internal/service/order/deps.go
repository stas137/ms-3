package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/stas137/ms-3/order/internal/model"
)

type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

// OrderRepository определяет контракт для работы с хранилищем заказов
// Интерфейс определен по месту использования - в пакете service
type OrderRepository interface {
	Create(ctx context.Context, order model.Order, items []model.OrderItem) error
	Get(ctx context.Context, orderUUID uuid.UUID) (model.Order, error)
	Update(ctx context.Context, order model.Order) error
}

type InventoryClient interface {
	ListParts(ctx context.Context, uuids []string) ([]model.Part, error)
	ValidateCompatibility(ctx context.Context, hullUUID, engineUUID, shieldUUID, weaponUUID string) error
	ReserveParts(ctx context.Context, uuids []string) error
	ReleaseParts(ctx context.Context, uuids []string) error
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID string, method model.PaymentMethod) (string, error)
}
