package order

import (
	"context"

	"github.com/stas137/ms-3/order/internal/model"
	orderv1 "github.com/stas137/ms-3/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(
	ctx context.Context,
	req *orderv1.PayOrderRequest,
	params orderv1.PayOrderParams,
) (orderv1.PayOrderRes, error) {
	transactionUUID, err := a.orderService.Pay(ctx, params.OrderUUID, model.PaymentMethod(req.PaymentMethod))
	if err != nil {
		return &orderv1.PayOrderResponse{}, err
	}
	return &orderv1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}
