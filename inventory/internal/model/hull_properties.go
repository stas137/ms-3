package model

import (
	"fmt"

	errs "github.com/stas137/ms-3/inventory/internal/errors"
)

type HullProperties struct {
	strength int
}

func (h *HullProperties) Strength() int {
	return h.strength
}

func (h *HullProperties) CanSupport(e *EngineProperties) bool {
	return h.Strength() >= e.RequiredStrength()
}

// NewHullProperties создает свойства корпуса.
// Прочность должна быть в диапазоне 30-200
func NewHullProperties(strength int) (PartProperties, error) {
	if strength < 30 || strength > 200 {
		return PartProperties{}, fmt.Errorf("прочность корпуса должна быть от 30 до 200, получено %d, %w", strength, errs.ErrInvalidProperties)
	}

	return PartProperties{
		hull: &HullProperties{
			strength: strength,
		},
	}, nil
}
