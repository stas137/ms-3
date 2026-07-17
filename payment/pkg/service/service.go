package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	paymentv1 "github.com/stas137/ms-3/shared/pkg/proto/payment/v1"
)

// server реализует gRPC сервис оплаты
type server struct {
	paymentv1.UnimplementedPaymentServiceServer
}

// NewServer создаёт новый экземпляр сервера оплаты
func NewServer() *server {
	return &server{}
}

// PayOrder обрабатывает оплату заказа
func (s *server) PayOrder(
	ctx context.Context,
	req *paymentv1.PayOrderRequest,
) (*paymentv1.PayOrderResponse, error) {
	// Реализовать метод
	// 1. Проверить, что order_uuid не пустой → INVALID_ARGUMENT
	// 2. Проверить, что payment_method != UNSPECIFIED → INVALID_ARGUMENT
	// 3. Проверить формат UUID → INVALID_ARGUMENT
	// 4. Сгенерировать transaction_uuid (UUID v4)
	// 5. Вывести в лог: "оплата прошла успешно, order_uuid: X, transaction_uuid: Y"
	// 6. Вернуть transaction_uuid

	if req.GetOrderUuid() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "неверный формат uuid: %s", req.GetOrderUuid())
	}
	if req.GetPaymentMethod() == paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED {
		return nil, status.Errorf(codes.InvalidArgument, "неверный метод оплаты: %s", req.GetPaymentMethod())
	}
	if _, err := uuid.Parse(req.GetOrderUuid()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "неверный формат uuid: %s", req.GetOrderUuid())
	}

	transaction_uuid := uuid.New()

	slog.Info(
		"оплата прошла успешно",
		"order_uuid", req.GetOrderUuid(),
		"transaction_uuid", transaction_uuid.String(),
	)

	return &paymentv1.PayOrderResponse{
		TransactionUuid: transaction_uuid.String(),
	}, nil

	// return nil, status.Error(codes.Unimplemented, "метод PayOrder не реализован")
}
