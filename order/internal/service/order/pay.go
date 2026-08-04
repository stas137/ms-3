package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	errs "github.com/stas137/ms-3/order/internal/errors"
	"github.com/stas137/ms-3/order/internal/model"
)

func (s *service) Pay(
	ctx context.Context,
	orderUUID uuid.UUID,
	method model.PaymentMethod,
) (uuid.UUID, error) {
	order, err := s.orderRepo.Get(ctx, orderUUID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("оплатить заказ: %w", err)
	}
	if order.Status == model.OrderStatusPaid {
		return uuid.Nil, errs.ErrOrderAlreadyPaid
	}
	if order.Status == model.OrderStatusCancelled {
		return uuid.Nil, errs.ErrOrderCancelled
	}

	transactionUUID, err := s.paymentClient.PayOrder(ctx, orderUUID.String(), method)
	if err != nil {
		return uuid.Nil, fmt.Errorf("оплатить заказ: %w", err)
	}

	parsedUUID, err := uuid.Parse(transactionUUID)
	if err != nil {
		return uuid.Nil, errs.ErrInvalidUUID
	}

	modelOrder := model.Order{
		UUID:            orderUUID,
		Items:           order.Items,
		TransactionUUID: (&parsedUUID),
		PaymentMethod:   &method,
		Status:          model.OrderStatusPaid,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       new(time.Now()),
	}

	err = s.orderRepo.Update(ctx, modelOrder)
	if err != nil {
		return uuid.Nil, fmt.Errorf("оплатить заказ: %w", err)
	}

	return parsedUUID, nil
}
