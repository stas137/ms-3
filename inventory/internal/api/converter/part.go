package converter

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	errs "github.com/stas137/ms-3/inventory/internal/errors"
	"github.com/stas137/ms-3/inventory/internal/model"
	inventoryv1 "github.com/stas137/ms-3/shared/pkg/proto/inventory/v1"
)

func PartToDTO(part model.Part) *inventoryv1.Part {
	inventoryv1PartTypeName := "PART_TYPE_" + string(part.PartType())
	inventoryv1PartType := inventoryv1.PartType_value[inventoryv1PartTypeName]

	return &inventoryv1.Part{
		Uuid:          part.UUID().String(),
		Name:          part.Name(),
		Description:   part.Description(),
		Price:         part.Price(),
		PartType:      inventoryv1.PartType(inventoryv1PartType),
		StockQuantity: int64(part.StockQuantity()),
		CreatedAt:     timestamppb.New(part.CreatedAt()),
	}
}

func PartsToDTOs(parts []model.Part) [](*inventoryv1.Part) {
	res := make([](*inventoryv1.Part), len(parts))
	for idx, part := range parts {
		res[idx] = PartToDTO(part)
	}
	return res
}

func PartTypeToModelPartType(partType inventoryv1.PartType) model.PartType {
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
		return model.PartTypeUnspecified
	}
}

func StringToUUID(s string) (uuid.UUID, error) {
	parsedUUID, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, errs.ErrInvalidUUID
	}
	return parsedUUID, nil
}

func StringsToUUIDs(s []string) ([]uuid.UUID, error) {
	var partsUUID []uuid.UUID
	for _, partUUID := range s {
		if partUUID != "" {
			parsedUUID, err := uuid.Parse(partUUID)
			if err != nil {
				return nil, errs.ErrInvalidUUID
			}
			partsUUID = append(partsUUID, parsedUUID)
		}
	}
	return partsUUID, nil
}
