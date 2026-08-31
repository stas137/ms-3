package model

import (
	"fmt"

	errs "github.com/stas137/ms-3/inventory/internal/errors"
)

type EngineClass string

type EngineProperties struct {
	class             EngineClass
	required_strength int
}

func (e *EngineProperties) Class() EngineClass {
	return e.class
}

func (e *EngineProperties) RequiredStrength() int {
	return e.required_strength
}

func NewEngineProperties(class EngineClass, required_strength int) (PartProperties, error) {
	if class != "A" && class != "B" && class != "C" {
		return PartProperties{}, fmt.Errorf("класс двигателя должен быть А, В или С, получено %s, %w", class, errs.ErrInvalidProperties)
	}
	if required_strength < 30 {
		return PartProperties{}, fmt.Errorf("требуемая прочность двигателя должна быть больше 30, получено %d, %w", required_strength, errs.ErrInvalidProperties)
	}

	if class == "C" && required_strength == 30 {
		return PartProperties{
			engine: &EngineProperties{
				class:             EngineClass(class),
				required_strength: required_strength,
			},
		}, nil
	}
	if class == "B" && required_strength == 70 {
		return PartProperties{
			engine: &EngineProperties{
				class:             EngineClass(class),
				required_strength: required_strength,
			},
		}, nil
	}
	if class == "A" && required_strength == 100 {
		return PartProperties{
			engine: &EngineProperties{
				class:             EngineClass(class),
				required_strength: required_strength,
			},
		}, nil
	}

	return PartProperties{}, fmt.Errorf("требуемая прочность двигателя и корпуса не совпадает, получено %s %d, %w", class, required_strength, errs.ErrInvalidProperties)
}
