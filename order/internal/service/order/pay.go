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
	var transactionUUID uuid.UUID

	err := s.txManager.Do(ctx, func(txCtx context.Context) error {
		order, err := s.orderRepo.Get(ctx, orderUUID)
		if err != nil {
			return fmt.Errorf("получить заказ: %w", err)
		}

		if order.Status != model.OrderStatusPendingPayment {

			if order.Status == model.OrderStatusPaid {
				return errs.ErrOrderAlreadyPaid
			}
			return errs.ErrOrderCancelled
		}

		tempUUID, err := s.paymentClient.PayOrder(ctx, orderUUID.String(), method)
		if err != nil {
			return fmt.Errorf("оплатить заказ: %w", err)
		}

		transactionUUID, err = uuid.Parse(tempUUID)
		if err != nil {
			return errs.ErrInvalidUUID
		}

		modelOrder := model.Order{
			UUID:            orderUUID,
			Items:           order.Items,
			TransactionUUID: (&transactionUUID),
			PaymentMethod:   &method,
			Status:          model.OrderStatusPaid,
			CreatedAt:       order.CreatedAt,
			UpdatedAt:       new(time.Now()),
		}

		err = s.orderRepo.Update(ctx, modelOrder)
		if err != nil {
			return fmt.Errorf("оплатить заказ: %w", err)
		}

		return nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	return transactionUUID, nil
}
