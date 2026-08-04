package part

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"

	errs "github.com/stas137/ms-3/inventory/internal/errors"
	"github.com/stas137/ms-3/inventory/internal/model"
	"github.com/stas137/ms-3/inventory/internal/repository/converter"
	"github.com/stas137/ms-3/inventory/internal/repository/record"
)

type repository struct {
	mu    sync.RWMutex
	parts map[uuid.UUID]record.Part
}

func NewRepository() *repository {
	return &repository{
		parts: map[uuid.UUID]record.Part{
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440001",
				Name:          "Алюминиевый корпус",
				Description:   "Лёгкий корпус для небольших кораблей",
				Price:         500000, // 5000₽
				PartType:      "HULL",
				StockQuantity: 10,
				CreatedAt:     time.Now(),
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440002",
				Name:          "Титановый корпус",
				Description:   "Прочный корпус для средних кораблей",
				Price:         1500000, // 15000₽
				PartType:      "HULL",
				StockQuantity: 5,
				CreatedAt:     time.Now(),
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440003"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440003",
				Name:          "Ионный двигатель C",
				Description:   "Базовый ионный двигатель класса C",
				Price:         300000, // 3000₽
				PartType:      "ENGINE",
				StockQuantity: 8,
				CreatedAt:     time.Now(),
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440004"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440004",
				Name:          "Ионный двигатель B",
				Description:   "Улучшенный ионный двигатель класса B",
				Price:         800000, // 8000₽
				PartType:      "ENGINE",
				StockQuantity: 3,
				CreatedAt:     time.Now(),
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440005"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440005",
				Name:          "Энергетический щит",
				Description:   "Стандартный энергетический щит",
				Price:         400000, // 4000₽
				PartType:      "SHIELD",
				StockQuantity: 6,
				CreatedAt:     time.Now(),
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440006"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440006",
				Name:          "Лазерная пушка",
				Description:   "Точная лазерная пушка",
				Price:         250000, // 2500₽
				PartType:      "WEAPON",
				StockQuantity: 7,
				CreatedAt:     time.Now(),
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440007"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440007",
				Name:          "Плазменный корпус",
				Description:   "Экспериментальный корпус (нет на складе)",
				Price:         2000000, // 20000₽
				PartType:      "HULL",
				StockQuantity: 0,
				CreatedAt:     time.Now(),
			},
		},
	}
}

func (r *repository) Get(_ context.Context, partUUID uuid.UUID) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	part, ok := r.parts[partUUID]
	if !ok {
		return model.Part{}, errs.ErrPartNotFound
	}

	modelPart, err := converter.PartToModelPart(part)
	if err != nil {
		return model.Part{}, err
	}
	return modelPart, nil
}

func (r *repository) List(_ context.Context, partsUUID []uuid.UUID, partType model.PartType) ([]model.Part, error) {
	var res []model.Part

	if len(partsUUID) != 0 {
		for _, id := range partsUUID {
			if detail, ok := r.parts[id]; ok {
				modelPart, err := converter.PartToModelPart(detail)
				if err != nil {
					return nil, err
				}
				res = append(res, modelPart)
			} else {
				return nil, errs.ErrPartNotFound
			}
		}
		return res, nil
	}

	for _, detail := range r.parts {
		if detail.PartType == string(partType) || partType == model.PartTypeUnspecified {
			modelPart, err := converter.PartToModelPart(detail)
			if err != nil {
				return nil, err
			}
			res = append(res, modelPart)
		}
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})
	return res, nil
}
