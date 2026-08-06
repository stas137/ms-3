package payment

import (
	"context"

	"github.com/google/uuid"

	"github.com/stas137/ms-3/payment/internal/model"
)

type PaymentRepository interface {
	Pay(ctx context.Context, orderUUID uuid.UUID, paymentMethod model.PaymentMethod) (uuid.UUID, error)
}
