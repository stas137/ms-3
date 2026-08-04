package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/stas137/ms-3/order/internal/api/converter"
	errs "github.com/stas137/ms-3/order/internal/errors"
	orderv1 "github.com/stas137/ms-3/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrder(
	ctx context.Context,
	params orderv1.GetOrderParams,
) (
	orderv1.GetOrderRes,
	error,
) {
	if params.OrderUUID == uuid.Nil {
		return nil, errs.ErrInvalidUUID
	}
	orderUUID, err := converter.ToGetInput(params.OrderUUID.String())
	if err != nil {
		return nil, err
	}
	order, err := a.orderService.Get(ctx, orderUUID)
	if err != nil {
		return nil, err
	}

	return converter.OrderToDTO(order), nil
}
