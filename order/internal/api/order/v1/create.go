package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/stas137/ms-3/order/internal/service/input"
	orderv1 "github.com/stas137/ms-3/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(
	ctx context.Context,
	req *orderv1.CreateOrderRequest,
) (
	orderv1.CreateOrderRes,
	error,
) {
	hullUUID := req.GetHullUUID()
	engineUUID := req.GetEngineUUID()

	var shieldUUID *uuid.UUID
	if value, ok := req.ShieldUUID.Get(); ok {
		shieldUUID = &value
	}
	var weaponUUID *uuid.UUID
	if value, ok := req.WeaponUUID.Get(); ok {
		weaponUUID = &value
	}

	order, err := a.orderService.Create(ctx, input.CreateOrderInput{
		HullUUID:   hullUUID,
		EngineUUID: engineUUID,
		ShieldUUID: shieldUUID,
		WeaponUUID: weaponUUID,
	})
	if err != nil {
		return &orderv1.CreateOrderResponse{}, err
	}

	return &orderv1.CreateOrderResponse{
		OrderUUID:  order.UUID,
		TotalPrice: order.TotalPrice(),
	}, nil
}
