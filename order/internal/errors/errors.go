package errs

import "errors"

var (
	// Ошибки заказов
	ErrOrderNotFound    = errors.New("заказ не найден")
	ErrOrderAlreadyPaid = errors.New("заказ уже оплачен")
	ErrOrderCancelled   = errors.New("заказ отменен")

	// Ошибки деталей
	ErrPartNotFound = errors.New("деталь не найдена")
	ErrOutOfStock   = errors.New("деталь отсутствует на складе")

	// Ошибки валидации
	ErrInvalidUUID          = errors.New("неверный формат UUID")
	ErrInvalidPaymentMethod = errors.New("неверный метод оплаты")

	// Детали ошибки
	ErrIncompatibleParts = errors.New("детали несовместимы")
	ErrPartTypeMismatch  = errors.New("тип детали не соответствует слоту корабля")
)
