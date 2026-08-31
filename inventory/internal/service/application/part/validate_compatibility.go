package part

import (
	"context"

	"github.com/google/uuid"

	errs "github.com/stas137/ms-3/inventory/internal/errors"
	"github.com/stas137/ms-3/inventory/internal/model"
)

func (s *service) ValidateCompatibility(ctx context.Context, hullUUID, engineUUID uuid.UUID, shieldUUID, weaponUUID *uuid.UUID) error {
	resolvedShipSlots, err := s.resolvedShipSlots(ctx, hullUUID, engineUUID, shieldUUID, weaponUUID)
	if err != nil {
		return err
	}

	return s.compatibilityChecker.Check(resolvedShipSlots)
}

func (s *service) resolvedShipSlots(ctx context.Context, hullUUID, engineUUID uuid.UUID, shieldUUID, weaponUUID *uuid.UUID) (model.ShipSlots, error) {
	uuids := []uuid.UUID{
		hullUUID,
		engineUUID,
	}
	if shieldUUID != nil {
		uuids = append(uuids, *shieldUUID)
	}
	if weaponUUID != nil {
		uuids = append(uuids, *weaponUUID)
	}

	parts, err := s.partRepository.List(ctx, uuids, model.PartTypeUnspecified)
	if err != nil {
		return model.ShipSlots{}, err
	}

	partsMap := make(map[uuid.UUID]struct{})
	resolvedShipSlot := make(map[model.PartType]model.Part)

	for _, part := range parts {
		if part.UUID() == hullUUID && part.PartType() != model.PartTypeHull {
			return model.ShipSlots{}, errs.ErrPartTypeMismatch
		}
		if part.UUID() == engineUUID && part.PartType() != model.PartTypeEngine {
			return model.ShipSlots{}, errs.ErrPartTypeMismatch
		}
		if shieldUUID != nil && part.UUID() == *shieldUUID && part.PartType() != model.PartTypeShield {
			return model.ShipSlots{}, errs.ErrPartTypeMismatch
		}
		if weaponUUID != nil && part.UUID() == *weaponUUID && part.PartType() != model.PartTypeWeapon {
			return model.ShipSlots{}, errs.ErrPartTypeMismatch
		}
		if _, ok := partsMap[part.UUID()]; ok {
			return model.ShipSlots{}, errs.ErrPartTypeMismatch
		}
		partsMap[part.UUID()] = struct{}{}
		resolvedShipSlot[part.PartType()] = part
	}

	var shield *model.Part
	if _, ok := resolvedShipSlot[model.PartTypeShield]; ok {
		shield = new(resolvedShipSlot[model.PartTypeShield])
	}
	var weapon *model.Part
	if _, ok := resolvedShipSlot[model.PartTypeWeapon]; ok {
		weapon = new(resolvedShipSlot[model.PartTypeWeapon])
	}

	return model.ShipSlots{
		Hull:   resolvedShipSlot[model.PartTypeHull],
		Engine: resolvedShipSlot[model.PartTypeEngine],
		Shield: shield,
		Weapon: weapon,
	}, nil
}
