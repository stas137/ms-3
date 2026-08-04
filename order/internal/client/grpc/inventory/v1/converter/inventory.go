package converter

import (
	"github.com/google/uuid"

	errs "github.com/stas137/ms-3/order/internal/errors"
	"github.com/stas137/ms-3/order/internal/model"
	inventoryv1 "github.com/stas137/ms-3/shared/pkg/proto/inventory/v1"
)

func ProtoPartsToModelParts(parts []*inventoryv1.Part) ([]model.Part, error) {
	// func ProtoPartsToModelParts(parts []*inventoryv1.Part) []model.Part {
	modelParts := make([]model.Part, len(parts))
	for idx, part := range parts {

		partUUID, err := uuid.Parse(part.GetUuid())
		if err != nil {
			return nil, errs.ErrInvalidUUID
		}

		modelParts[idx] = model.Part{
			UUID:          partUUID,
			Name:          part.GetName(),
			PartType:      protoPartTypeToModelPartType(part.PartType),
			Price:         part.GetPrice(),
			StockQuantity: part.StockQuantity,
		}
	}
	return modelParts, nil
}

func protoPartTypeToModelPartType(partType inventoryv1.PartType) model.PartType {
	switch partType {
	case inventoryv1.PartType_PART_TYPE_HULL:
		return model.PartTypeHull
	case inventoryv1.PartType_PART_TYPE_ENGINE:
		return model.PartTypeEngine
	case inventoryv1.PartType_PART_TYPE_SHIELD:
		return model.PartTypeShield
	case inventoryv1.PartType_PART_TYPE_WEAPON:
		return model.PartTypeWeapon
	default:
		return model.PartType("")
	}
}
