package input

import (
	"github.com/google/uuid"

	"github.com/stas137/ms-3/payment/internal/model"
)

type PayOrderInput struct {
	OrderUUID     uuid.UUID
	PaymentMethod model.PaymentMethod
}
