package inventory

import (
	"context"

	"github.com/stas137/ms-3/inventory/internal/api/converter"
	inventoryv1 "github.com/stas137/ms-3/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *inventoryv1.GetPartRequest,
) (*inventoryv1.GetPartResponse, error) {
	parsedUUID, err := converter.StringToUUID(req.Uuid)
	if err != nil {
		return &inventoryv1.GetPartResponse{}, err
	}

	part, err := a.partService.Get(ctx, parsedUUID)
	if err != nil {
		return &inventoryv1.GetPartResponse{}, err
	}
	partDTO := converter.PartToDTO(part)
	return &inventoryv1.GetPartResponse{
		Part: partDTO,
	}, nil
}
