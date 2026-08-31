package part

import (
	"context"
	"errors"
	"sort"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	errs "github.com/stas137/ms-3/inventory/internal/errors"
	"github.com/stas137/ms-3/inventory/internal/model"
	"github.com/stas137/ms-3/inventory/internal/repository/converter"
	"github.com/stas137/ms-3/inventory/internal/repository/record"
)

type repository struct {
	invPool *pgxpool.Pool
}

func NewRepository(invPool *pgxpool.Pool) *repository {
	return &repository{
		invPool: invPool,
	}
}

func (r *repository) listByUUIDs(ctx context.Context, partsUUID []uuid.UUID) ([]model.Part, error) {
	var partDTO record.PartRecord
	var res []model.Part

	query := `SELECT 
				uuid, 
				name, 
				description, 
				part_type, 
				price, 
				stock_quantity,
				reserved,
				properties,
				created_at, 
				updated_at
			FROM parts 
			WHERE uuid = $1`

	for _, partUUID := range partsUUID {
		row := r.invPool.QueryRow(ctx, query, partUUID)
		err := row.Scan(
			&partDTO.UUID,
			&partDTO.Name,
			&partDTO.Description,
			&partDTO.PartType,
			&partDTO.Price,
			&partDTO.StockQuantity,
			&partDTO.Reserved,
			&partDTO.Properties,
			&partDTO.CreatedAt,
			&partDTO.UpdatedAt,
		)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, errs.ErrPartNotFound
			}
			return nil, err
		}
		modelPart, err := converter.PartRecordToModelPart(partDTO)
		if err != nil {
			return nil, err
		}
		res = append(res, modelPart)
	}
	return res, nil
}

func (r *repository) listByPartTypeUnspecified(ctx context.Context) ([]model.Part, error) {
	var partDTO record.PartRecord
	var res []model.Part

	query := `SELECT 
				uuid, 
				name, 
				description, 
				part_type, 
				price, 
				stock_quantity,
				reserved,
				properties,
				created_at, 
				updated_at 
			FROM parts`

	rows, err := r.invPool.Query(
		ctx,
		query,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&partDTO.UUID,
			&partDTO.Name,
			&partDTO.Description,
			&partDTO.PartType,
			&partDTO.Price,
			&partDTO.StockQuantity,
			&partDTO.Reserved,
			&partDTO.Properties,
			&partDTO.CreatedAt,
			&partDTO.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		modelPart, err := converter.PartRecordToModelPart(partDTO)
		if err != nil {
			return nil, err
		}
		res = append(res, modelPart)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name() < res[j].Name()
	})
	return res, nil
}

func (r *repository) listByPartType(ctx context.Context, partType model.PartType) ([]model.Part, error) {
	var partDTO record.PartRecord
	var res []model.Part

	query := `
		SELECT 
			uuid, 
			name, 
			description, 
			part_type, 
			price, 
			stock_quantity,
			reserved,
			properties,
			created_at, 
			updated_at 
		FROM parts 
		WHERE part_type = $1`

	rows, err := r.invPool.Query(
		ctx,
		query,
		partType,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&partDTO.UUID,
			&partDTO.Name,
			&partDTO.Description,
			&partDTO.PartType,
			&partDTO.Price,
			&partDTO.StockQuantity,
			&partDTO.Reserved,
			&partDTO.Properties,
			&partDTO.CreatedAt,
			&partDTO.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		modelPart, err := converter.PartRecordToModelPart(partDTO)
		if err != nil {
			return nil, err
		}
		res = append(res, modelPart)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name() < res[j].Name()
	})
	return res, nil
}

func (r *repository) Get(ctx context.Context, partUUID uuid.UUID) (model.Part, error) {
	var partDTO record.PartRecord
	query := `
		SELECT 
			uuid, 
			name, 
			description, 
			part_type, 
			price, 
			stock_quantity,
			reserved,
			properties,
			created_at, 
			updated_at 
		FROM parts 
		WHERE uuid = $1`

	err := r.invPool.QueryRow(ctx, query, partUUID).Scan(
		&partDTO.UUID,
		&partDTO.Name,
		&partDTO.Description,
		&partDTO.PartType,
		&partDTO.Price,
		&partDTO.StockQuantity,
		&partDTO.Reserved,
		&partDTO.Properties,
		&partDTO.CreatedAt,
		&partDTO.UpdatedAt,
	)
	if err != nil {
		return model.Part{}, errs.ErrPartNotFound
	}
	partModel, err := converter.PartRecordToModelPart(partDTO)
	if err != nil {
		return model.Part{}, err
	}
	return partModel, nil
}

func (r *repository) List(ctx context.Context, partsUUID []uuid.UUID, partType model.PartType) ([]model.Part, error) {
	if len(partsUUID) != 0 {
		return r.listByUUIDs(ctx, partsUUID)
	}
	if partType == model.PartTypeUnspecified {
		return r.listByPartTypeUnspecified(ctx)
	}
	return r.listByPartType(ctx, partType)
}

func (r *repository) UpdateReservedBatch(ctx context.Context, parts []model.Part) error {
	query := `
	UPDATE parts AS p 
	SET reserved = batch.reserved,
		updated_at = NOW()
	FROM unnest($1::uuid[], $2::int[]) as batch(uuid, reserved)
	WHERE p.uuid = batch.uuid
	`

	uuids := make([]uuid.UUID, len(parts))
	reservedVals := make([]int, len(parts))

	for i, part := range parts {
		uuids[i] = part.UUID()
		reservedVals[i] = part.Reserved()
	}

	_, err := r.invPool.Exec(ctx, query, uuids, reservedVals)
	if err != nil {
		return err
	}

	return nil
}
