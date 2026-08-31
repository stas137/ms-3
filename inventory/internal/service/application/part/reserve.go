package part

import (
	"context"

	"github.com/google/uuid"

	"github.com/stas137/ms-3/inventory/internal/model"
)

func (s *service) ReserveParts(ctx context.Context, uuids []uuid.UUID) error {
	err := s.txManager.Do(ctx, func(txCtx context.Context) error {
		parts, err := s.partRepository.List(ctx, uuids, model.PartTypeUnspecified)
		if err != nil {
			return err
		}
		updatedParts := make([]model.Part, len(parts))
		for idx, part := range parts {
			err := part.Reserve(1)
			if err != nil {
				return err
			}
			updatedParts[idx] = part
		}
		err = s.partRepository.UpdateReservedBatch(ctx, updatedParts)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
