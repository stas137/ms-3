package inventory

import (
	"context"

	"github.com/google/uuid"

	"github.com/stas137/ms-3/inventory/internal/model"
	"github.com/stas137/ms-3/inventory/internal/service/input"
)

type PartService interface {
	Get(ctx context.Context, partUUID uuid.UUID) (model.Part, error)
	List(ctx context.Context, filter input.PartFilter) ([]model.Part, error)
	ValidateCompatibility(ctx context.Context, hullUUID, engineUUID uuid.UUID, shieldUUID, weaponUUID *uuid.UUID) error
	ReserveParts(ctx context.Context, uuids []uuid.UUID) error
	ReleaseParts(ctx context.Context, uuids []uuid.UUID) error
}
