package model

import "github.com/google/uuid"

type PaymentMethod string

const (
	PaymentMethodUnspecified   PaymentMethod = "UNSPECIFIED"
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

type Payment struct {
	OrderUUID     uuid.UUID
	PaymentMethod PaymentMethod
}

func (pm PaymentMethod) IsValid() bool {
	switch pm {
	case
		PaymentMethodCard,
		PaymentMethodSBP,
		PaymentMethodCreditCard,
		PaymentMethodInvestorMoney:
		return true
	default:
		return false
	}
}
