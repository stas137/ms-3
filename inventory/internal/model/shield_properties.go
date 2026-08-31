package model

import (
	"fmt"

	errs "github.com/stas137/ms-3/inventory/internal/errors"
)

type ShieldType string

type ShieldProperties struct {
	shieldType ShieldType
}

func (s *ShieldProperties) ShieldType() ShieldType {
	return s.shieldType
}

func (s *ShieldProperties) ConflictWith(w *WeaponProperties) bool {
	if s.ShieldType() == "plasma" && w.WeaponType() == "laser" {
		return true
	}
	if s.ShieldType() == "energy" && w.WeaponType() == "laser" {
		return false
	}
	if s.ShieldType() == "plasma" && w.WeaponType() == "missier" {
		return false
	}
	if s.ShieldType() == "energy" && w.WeaponType() == "missier" {
		return false
	}
	return true
}

func NewShieldProperties(shieldType ShieldType) (PartProperties, error) {
	if shieldType != "energy" && shieldType != "plasma" {
		return PartProperties{}, fmt.Errorf("тип щита должен быть energy или plasma, получено %s: %w", shieldType, errs.ErrInvalidProperties)
	}

	return PartProperties{
		shield: &ShieldProperties{
			shieldType: ShieldType(shieldType),
		},
	}, nil
}
