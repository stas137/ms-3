package order

import (
	"context"

	orderv1 "github.com/stas137/ms-3/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(
	ctx context.Context,
	params orderv1.CancelOrderParams,
) (orderv1.CancelOrderRes, error) {
	err := a.orderService.Cancel(ctx, params.OrderUUID)
	if err != nil {
		return &orderv1.CancelOrderResponse{}, err // nil, err
	}
	return &orderv1.CancelOrderResponse{}, nil // nil, nil
}
