package inventory

import (
	"context"
	"fmt"

	"github.com/stas137/ms-3/inventory/internal/api/converter"
	"github.com/stas137/ms-3/inventory/internal/service/input"
	inventoryv1 "github.com/stas137/ms-3/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryv1.ListPartsRequest) (
	*inventoryv1.ListPartsResponse,
	error,
) {
	partsUUIDs, err := converter.StringsToUUIDs(req.GetUuids())
	if err != nil {
		return &inventoryv1.ListPartsResponse{}, err
	}

	partType := converter.PartTypeToModelPartType(req.GetPartType())

	parts, err := a.partService.List(ctx, input.PartFilter{
		UUIDs:    partsUUIDs,
		PartType: partType,
	})
	if err != nil {
		return &inventoryv1.ListPartsResponse{}, fmt.Errorf("получить детали: %w", err)
	}

	partsDTO := converter.PartsToDTOs(parts)

	return &inventoryv1.ListPartsResponse{
		Parts: partsDTO,
	}, nil
}
