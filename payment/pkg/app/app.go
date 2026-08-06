package app

import (
	"google.golang.org/grpc"

	paymentv1API "github.com/stas137/ms-3/payment/internal/api/payment/v1"
	"github.com/stas137/ms-3/payment/internal/interceptor"
	paymentService "github.com/stas137/ms-3/payment/internal/service/payment"
	paymentv1 "github.com/stas137/ms-3/shared/pkg/proto/payment/v1"
)

func RegisterServices(grpcServer *grpc.Server) {
	// paymentRepo := paymentRepository.NewRepository()
	paymentServ := paymentService.NewService()
	paymentApi := paymentv1API.NewApi(paymentServ)

	paymentv1.RegisterPaymentServiceServer(grpcServer, paymentApi)
}

func Interceptors() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.ErrorInterceptor),
	}
}
