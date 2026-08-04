package converter

import (
	"github.com/google/uuid"

	errs "github.com/stas137/ms-3/payment/internal/errors"
	"github.com/stas137/ms-3/payment/internal/model"
	paymentv1 "github.com/stas137/ms-3/shared/pkg/proto/payment/v1"
)

func PaymentTypeToModelPaymentType(paymentMethod paymentv1.PaymentMethod) model.PaymentMethod {
	switch paymentMethod {
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCard
	case paymentv1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethodSBP
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCreditCard
	case paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnspecified
	}
}

func StringToUUID(s string) (uuid.UUID, error) {
	parsedUUID, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, errs.ErrInvalidUUID
	}
	return parsedUUID, nil
}
