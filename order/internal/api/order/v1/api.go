package order

type api struct {
	orderService OrderService
}

func NewApi(orderService OrderService) *api {
	return &api{
		orderService: orderService,
	}
}
