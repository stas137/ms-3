package model

import (
	"time"

	"github.com/google/uuid"

	errs "github.com/stas137/ms-3/inventory/internal/errors"
)

type Part struct {
	uuid          uuid.UUID
	name          string
	description   string
	price         int64
	partType      PartType
	stockQuantity int
	reserved      int
	properties    PartProperties
	createdAt     time.Time
}

func RestorePart(partUUID uuid.UUID, name, description string, partType PartType, price int64, stockQuantity, reserved int, properties PartProperties, createdAt time.Time) Part {
	return Part{
		uuid:          partUUID,
		name:          name,
		description:   description,
		partType:      partType,
		price:         price,
		stockQuantity: stockQuantity,
		reserved:      reserved,
		properties:    properties,
		createdAt:     createdAt,
	}
}

func (p *Part) Reserve(quantity int) error {
	if quantity < 0 {
		return errs.ErrOutOfStock
	}
	if p.StockQuantity() < quantity {
		return errs.ErrOutOfStock
	}
	if p.reserved+quantity > p.StockQuantity() {
		return errs.ErrOutOfStock
	}
	p.reserved += quantity

	return nil
}

func (p *Part) Release(quantity int) error {
	if p.Reserved() == 0 {
		return errs.ErrNothingToRelease
	}

	// if quantity < 0 {
	// 	return errs.ErrNothingToRelease
	// }
	// if quantity > p.Reserved() {
	// 	return errs.ErrNothingToRelease
	// }
	p.reserved -= quantity

	return nil
}

func (p *Part) UUID() uuid.UUID {
	return p.uuid
}

func (p *Part) Name() string {
	return p.name
}

func (p *Part) Description() string {
	return p.description
}

func (p *Part) PartType() PartType {
	return p.partType
}

func (p *Part) Price() int64 {
	return p.price
}

func (p *Part) StockQuantity() int {
	return p.stockQuantity
}

func (p *Part) Reserved() int {
	return p.reserved
}

func (p *Part) Properties() PartProperties {
	return p.properties
}

func (p *Part) CreatedAt() time.Time {
	return p.createdAt
}
