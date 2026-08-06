package errs

import "errors"

var (
	ErrInvalidUUID          = errors.New("невеный формат UUID")
	ErrInvalidPaymentMethod = errors.New("невеный метод оплаты")
)
