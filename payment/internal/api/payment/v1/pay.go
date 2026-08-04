package payment

import (
	"context"
	"fmt"

	"github.com/stas137/ms-3/payment/internal/api/converter"
	errs "github.com/stas137/ms-3/payment/internal/errors"
	"github.com/stas137/ms-3/payment/internal/service/input"
	paymentv1 "github.com/stas137/ms-3/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	orderUUID, err := converter.StringToUUID(req.GetOrderUuid())
	if err != nil {
		return nil, err
	}

	paymentMethod := converter.PaymentTypeToModelPaymentType(req.PaymentMethod)
	if !paymentMethod.IsValid() {
		return nil, errs.ErrInvalidPaymentMethod
	}

	transactionUUID, err := a.paymentService.Pay(ctx, input.PayOrderInput{
		OrderUUID:     orderUUID,
		PaymentMethod: paymentMethod,
	})
	if err != nil {
		return nil, fmt.Errorf("оплатить заказ: %w", err)
	}

	return &paymentv1.PayOrderResponse{
		TransactionUuid: transactionUUID.String(),
	}, nil
}
