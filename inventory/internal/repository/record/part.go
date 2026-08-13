package record

import (
	"time"
)

type Part struct {
	UUID          string
	Name          string
	Description   string
	PartType      string
	Price         int64
	StockQuantity int64
	CreatedAt     time.Time
	UpdatedAt     *time.Time
}
