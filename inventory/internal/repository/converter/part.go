package converter

import (
	"github.com/google/uuid"

	errs "github.com/stas137/ms-3/inventory/internal/errors"
	"github.com/stas137/ms-3/inventory/internal/model"
	"github.com/stas137/ms-3/inventory/internal/repository/record"
)

func PartToModelPart(part record.Part) (model.Part, error) {
	partUUID, err := uuid.Parse(part.UUID)
	if err != nil {
		return model.Part{}, errs.ErrInvalidUUID
	}

	return model.Part{
		UUID:          partUUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		PartType:      model.PartType(part.PartType),
		StockQuantity: part.StockQuantity,
		CreatedAt:     part.CreatedAt,
	}, nil
}
