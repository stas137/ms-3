package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

// PaymentMethod — способ оплаты заказа
type PaymentMethod string

const (
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

// Order представляет заказ на постройку космического корабля
type Order struct {
	UUID            uuid.UUID
	Items           []OrderItem
	TransactionUUID *uuid.UUID
	PaymentMethod   *PaymentMethod
	Status          OrderStatus
	CreatedAt       time.Time
	UpdatedAt       *time.Time
	DeletedAt       *time.Time
}

func (o Order) TotalPrice() int64 {
	var total int64
	for _, item := range o.Items {
		total += item.Price
	}
	return total
}

type OrderItem struct {
	PartUUID uuid.UUID
	PartType PartType
	Price    int64
}

func (o Order) GetEngineUUID() *uuid.UUID {
	for _, item := range o.Items {
		if item.PartType == PartTypeEngine {
			return &item.PartUUID
		}
	}
	return nil
}

func (o Order) GetHullUUID() *uuid.UUID {
	for _, item := range o.Items {
		if item.PartType == PartTypeHull {
			return &item.PartUUID
		}
	}
	return nil
}

func (o Order) GetShieldUUID() *uuid.UUID {
	for _, item := range o.Items {
		if item.PartType == PartTypeShield {
			return &item.PartUUID
		}
	}
	return nil
}

func (o Order) GetWeaponUUID() *uuid.UUID {
	for _, item := range o.Items {
		if item.PartType == PartTypeWeapon {
			return &item.PartUUID
		}
	}
	return nil
}
