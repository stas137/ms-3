package converter

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	errs "github.com/stas137/ms-3/inventory/internal/errors"
	"github.com/stas137/ms-3/inventory/internal/model"
	"github.com/stas137/ms-3/inventory/internal/repository/record"
)

func PartRecordToModelPart(rec record.PartRecord) (model.Part, error) {
	var propsRec record.PartPropertiesRecord
	if err := json.Unmarshal(rec.Properties, &propsRec); err != nil {
		return model.Part{}, fmt.Errorf("десереализовать свойства: %w", err)
	}

	partUUID, err := uuid.Parse(rec.UUID)
	if err != nil {
		return model.Part{}, errs.ErrInvalidUUID
	}

	props, err := partPropertiesFromRecord(propsRec)
	if err != nil {
		return model.Part{}, fmt.Errorf("коневертировать свойства: %w", err)
	}

	partType, err := model.NewPartType(rec.PartType)
	if err != nil {
		return model.Part{}, fmt.Errorf("конвертировать тип детали: %w", err)
	}

	return model.RestorePart(
		partUUID,
		rec.Name,
		rec.Description,
		partType,
		rec.Price,
		rec.StockQuantity,
		rec.Reserved,
		props,
		rec.CreatedAt,
	), nil
}

func partPropertiesFromRecord(rec record.PartPropertiesRecord) (model.PartProperties, error) {
	switch {
	case rec.Hull != nil:
		return model.NewHullProperties(rec.Hull.Strength)
	case rec.Engine != nil:
		return model.NewEngineProperties(model.EngineClass(rec.Engine.Class), rec.Engine.RequiredStrength)
	case rec.Shield != nil:
		return model.NewShieldProperties(model.ShieldType(rec.Shield.ShieldType))
	case rec.Weapon != nil:
		return model.NewWeaponProperties(model.WeaponType(rec.Weapon.WeaponType))
	default:
		return model.PartProperties{}, nil
	}
}
