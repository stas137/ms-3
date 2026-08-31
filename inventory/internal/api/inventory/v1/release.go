package inventory

import (
	"context"

	"github.com/stas137/ms-3/inventory/internal/api/converter"
	inventoryv1 "github.com/stas137/ms-3/shared/pkg/proto/inventory/v1"
)

func (a *api) ReleaseParts(ctx context.Context, req *inventoryv1.ReleasePartsRequest) (*inventoryv1.ReleasePartsResponse, error) {
	parsedUUIDs, err := converter.StringsToUUIDs(req.GetUuids())
	if err != nil {
		return &inventoryv1.ReleasePartsResponse{}, err
	}
	err = a.partService.ReleaseParts(ctx, parsedUUIDs)
	if err != nil {
		return &inventoryv1.ReleasePartsResponse{}, err
	}
	return &inventoryv1.ReleasePartsResponse{}, nil
}
