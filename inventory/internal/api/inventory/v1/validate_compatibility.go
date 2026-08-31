package inventory

import (
	"context"

	"github.com/google/uuid"

	"github.com/stas137/ms-3/inventory/internal/api/converter"
	inventoryv1 "github.com/stas137/ms-3/shared/pkg/proto/inventory/v1"
)

func (a *api) ValidateCompatibility(ctx context.Context, req *inventoryv1.ValidateCompatibilityRequest) (
	*inventoryv1.ValidateCompatibilityResponse, error,
) {
	hullUUID, err := converter.StringToUUID(req.HullUuid)
	if err != nil {
		return &inventoryv1.ValidateCompatibilityResponse{}, err
	}
	engineUUID, err := converter.StringToUUID(req.EngineUuid)
	if err != nil {
		return &inventoryv1.ValidateCompatibilityResponse{}, err
	}

	var shieldUUID *uuid.UUID
	var weaponUUID *uuid.UUID

	if req.ShieldUuid != "" {
		tempUUID, err := converter.StringToUUID(req.ShieldUuid)
		if err != nil {
			return &inventoryv1.ValidateCompatibilityResponse{}, err
		}
		shieldUUID = &tempUUID
	}
	if req.WeaponUuid != "" {
		tempUUID, err := converter.StringToUUID(req.WeaponUuid)
		if err != nil {
			return &inventoryv1.ValidateCompatibilityResponse{}, err
		}
		weaponUUID = &tempUUID
	}

	err = a.partService.ValidateCompatibility(ctx, hullUUID, engineUUID, shieldUUID, weaponUUID)
	if err != nil {
		return &inventoryv1.ValidateCompatibilityResponse{}, err
	}

	return &inventoryv1.ValidateCompatibilityResponse{}, nil
}
