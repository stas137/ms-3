package domain

import (
	errs "github.com/stas137/ms-3/inventory/internal/errors"
	"github.com/stas137/ms-3/inventory/internal/model"
)

type compatibilityChecker struct{}

func NewCompatibilityChecker() *compatibilityChecker {
	return &compatibilityChecker{}
}

// Check проверяет бизнес-правила совместимости для набора деталей
func (c *compatibilityChecker) Check(slots model.ShipSlots) error {
	hullProps := slots.Hull.Properties()
	engineProps := slots.Engine.Properties()

	canSupport := hullProps.Hull().CanSupport(engineProps.Engine())
	if !canSupport {
		return errs.ErrIncompatibleParts
	}

	var shieldProps model.PartProperties
	var weaponProps model.PartProperties

	if slots.Shield != nil && slots.Weapon != nil {
		shieldProps = slots.Shield.Properties()
		weaponProps = slots.Weapon.Properties()

		conflictWith := shieldProps.Shield().ConflictWith(weaponProps.Weapon())
		if conflictWith {
			return errs.ErrIncompatibleParts
		}
	}

	return nil
}
