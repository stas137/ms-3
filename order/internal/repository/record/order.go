package record

import (
	"time"

	"github.com/google/uuid"
)

// Order представляет заказ на постройку космического корабля
type Order struct {
	OrderUUID       uuid.UUID
	EngineUUID      uuid.UUID
	HullUUID        uuid.UUID
	ShieldUUID      *uuid.UUID // опциональный
	WeaponUUID      *uuid.UUID // опциональный
	EnginePrice     int64
	HullPrice       int64
	ShieldPrice     *int64
	WeaponPrice     *int64
	TotalPrice      int64 // в копейках
	TransactionUUID *uuid.UUID
	PaymentMethod   *string
	Status          string
	CreatedAt       time.Time
	UpdatedAt       *time.Time
	DeletedAt       *time.Time
}
