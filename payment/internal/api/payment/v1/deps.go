package payment

import (
	"context"

	"github.com/google/uuid"

	"github.com/stas137/ms-3/payment/internal/service/input"
)

type PaymentService interface {
	Pay(ctx context.Context, in input.PayOrderInput) (uuid.UUID, error)
}
