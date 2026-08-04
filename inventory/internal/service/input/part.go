package input

import (
	"github.com/google/uuid"

	"github.com/stas137/ms-3/inventory/internal/model"
)

type PartFilter struct {
	// UUIDs - если не пустой, возвращается только эти детали (приоритет)
	UUIDs []uuid.UUID
	// PartType - фильтр по типу (игнорируется если UUIDs заполнены)
	PartType model.PartType
}
