package order

import orderv1 "github.com/stas137/ms-3/shared/pkg/openapi/order/v1"

type api struct {
	orderv1.UnimplementedHandler
	orderService OrderService
}

func NewApi(orderService OrderService) *api {
	return &api{
		orderService: orderService,
	}
}
