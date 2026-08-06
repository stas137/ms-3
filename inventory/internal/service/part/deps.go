package part

import (
	"context"

	"github.com/google/uuid"

	"github.com/stas137/ms-3/inventory/internal/model"
)

type PartRepository interface {
	Get(ctx context.Context, partUUID uuid.UUID) (model.Part, error)
	List(ctx context.Context, partsUUID []uuid.UUID, partType model.PartType) ([]model.Part, error)
}
