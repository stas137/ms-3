package app

import (
	"context"

	partV1API "github.com/stas137/ms-3/payment/internal/api/payment/v1"
	paymentService "github.com/stas137/ms-3/payment/internal/service/payment"
	paymentv1 "github.com/stas137/ms-3/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	payServ          partV1API.PaymentService
	paymentV1Handler paymentv1.PaymentServiceServer
}

func (d *diContainer) PaymentService(ctx context.Context) partV1API.PaymentService {
	if d.payServ == nil {

		payServ := paymentService.NewService()
		d.payServ = payServ
	}
	return d.payServ
}

func (d *diContainer) PaymentV1API(ctx context.Context) paymentv1.PaymentServiceServer {
	if d.paymentV1Handler == nil {
		paymentV1Handler := partV1API.NewApi(d.PaymentService(ctx))
		d.paymentV1Handler = paymentV1Handler
	}
	return d.paymentV1Handler
}
