package payment

import (
	"context"

	"github.com/google/uuid"

	"github.com/stas137/ms-3/payment/internal/service/input"
)

func (s *service) Pay(ctx context.Context, in input.PayOrderInput) (uuid.UUID, error) {
	// orderUUID := in.OrderUUID
	// paymentMethod := in.PaymentMethod

	transactionUUID := uuid.New()

	return transactionUUID, nil
}
