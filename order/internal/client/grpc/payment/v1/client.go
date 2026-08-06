package v1

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errs "github.com/stas137/ms-3/order/internal/errors"
	"github.com/stas137/ms-3/order/internal/model"
	paymentv1 "github.com/stas137/ms-3/shared/pkg/proto/payment/v1"
)

type client struct {
	paymentClient paymentv1.PaymentServiceClient
}

func New(c paymentv1.PaymentServiceClient) *client {
	return &client{
		paymentClient: c,
	}
}

func (c client) PayOrder(
	ctx context.Context,
	orderUUID string,
	method model.PaymentMethod,
) (string, error) {
	paymentMethodName := "PAYMENT_METHOD_" + string(method)
	paymentv1PaymentMethodValue := paymentv1.PaymentMethod(paymentv1.PaymentMethod_value[paymentMethodName])

	resp, err := c.paymentClient.PayOrder(ctx, &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID,
		PaymentMethod: paymentv1PaymentMethodValue,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.NotFound {
				return "", errs.ErrOrderNotFound
			}
		}
		return "", fmt.Errorf("оплата заказа: %w", err)
	}

	transactionUUID := resp.GetTransactionUuid()

	return transactionUUID, nil
}
