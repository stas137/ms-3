package inventory

import (
	"context"

	"github.com/stas137/ms-3/inventory/internal/api/converter"
	inventoryv1 "github.com/stas137/ms-3/shared/pkg/proto/inventory/v1"
)

func (a *api) ReserveParts(ctx context.Context, req *inventoryv1.ReservePartsRequest) (*inventoryv1.ReservePartsResponse, error) {
	parsedUUIDs, err := converter.StringsToUUIDs(req.GetUuids())
	if err != nil {
		return &inventoryv1.ReservePartsResponse{}, err
	}

	err = a.partService.ReserveParts(ctx, parsedUUIDs)
	if err != nil {
		return &inventoryv1.ReservePartsResponse{}, err
	}
	return &inventoryv1.ReservePartsResponse{}, nil
}
