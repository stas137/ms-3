package part

import (
	"context"
	"fmt"

	"github.com/stas137/ms-3/inventory/internal/model"
	"github.com/stas137/ms-3/inventory/internal/service/input"
)

func (s *service) List(ctx context.Context, filter input.PartFilter,
) ([]model.Part, error) {
	parts, err := s.partRepository.List(ctx, filter.UUIDs, filter.PartType)
	if err != nil {
		return nil, fmt.Errorf("получить список деталей: %w", err)
	}
	return parts, nil
}
