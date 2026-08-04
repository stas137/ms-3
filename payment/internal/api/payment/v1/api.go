package payment

import paymentv1 "github.com/stas137/ms-3/shared/pkg/proto/payment/v1"

type api struct {
	paymentv1.UnimplementedPaymentServiceServer
	paymentService PaymentService
}

func NewApi(paymentService PaymentService) *api {
	return &api{
		paymentService: paymentService,
	}
}
