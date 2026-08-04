package converter

import (
	"github.com/google/uuid"

	errs "github.com/stas137/ms-3/order/internal/errors"
	"github.com/stas137/ms-3/order/internal/model"
	orderv1 "github.com/stas137/ms-3/shared/pkg/openapi/order/v1"
)

// DTO <-> Model + DTO -> Input

func OrderToDTO(order model.Order) orderv1.GetOrderRes {
	var shieldUUID orderv1.OptNilUUID
	if order.GetShieldUUID() != nil {
		shieldUUID = orderv1.NewOptNilUUID(*order.GetShieldUUID())
	}
	var weaponUUID orderv1.OptNilUUID
	if order.GetWeaponUUID() != nil {
		weaponUUID = orderv1.NewOptNilUUID(*order.GetWeaponUUID())
	}
	var transactionUUID orderv1.OptNilUUID
	if order.TransactionUUID != nil {
		transactionUUID = orderv1.NewOptNilUUID(*order.TransactionUUID)
	}
	var paymentMethod orderv1.OptNilPaymentMethod
	if order.PaymentMethod != nil {
		paymentMethod = orderv1.NewOptNilPaymentMethod(orderv1.PaymentMethod(*order.PaymentMethod))
	}

	return &orderv1.OrderDto{
		OrderUUID:       order.UUID,
		HullUUID:        *order.GetHullUUID(),
		EngineUUID:      *order.GetEngineUUID(),
		ShieldUUID:      shieldUUID,
		WeaponUUID:      weaponUUID,
		TotalPrice:      order.TotalPrice(),
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          orderv1.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
	}
}

// func (c *converter) ToGetInput(rawUUID string) (uuid.UUID, error) {
func ToGetInput(rawUUID string) (uuid.UUID, error) {
	parsed, err := uuid.Parse(rawUUID)
	if err != nil {
		return uuid.Nil, errs.ErrInvalidUUID
	}
	return parsed, nil
}
