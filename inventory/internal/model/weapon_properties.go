package model

import (
	"fmt"

	errs "github.com/stas137/ms-3/inventory/internal/errors"
)

type WeaponType string

type WeaponProperties struct {
	weaponType WeaponType
}

func (w *WeaponProperties) WeaponType() WeaponType {
	return w.weaponType
}

func NewWeaponProperties(weaponType WeaponType) (PartProperties, error) {
	if weaponType != "laser" && weaponType != "missile" {
		return PartProperties{}, fmt.Errorf("тип оружия должен быть laser или missile, получено %s: %w", weaponType, errs.ErrInvalidProperties)
	}

	return PartProperties{
		weapon: &WeaponProperties{
			weaponType: WeaponType(weaponType),
		},
	}, nil
}
