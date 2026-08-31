package part

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/stas137/ms-3/inventory/internal/model"
)

func (s *service) Get(ctx context.Context, partUUID uuid.UUID) (model.Part, error) {
	part, err := s.partRepository.Get(ctx, partUUID)
	if err != nil {
		return model.Part{}, fmt.Errorf("получить деталь: %w", err)
	}
	return part, nil
}
