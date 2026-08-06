package tests

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/stas137/ms-3/inventory/internal/api/converter"
	"github.com/stas137/ms-3/inventory/internal/api/inventory/v1"
	"github.com/stas137/ms-3/inventory/internal/api/inventory/v1/mocks"
	errs "github.com/stas137/ms-3/inventory/internal/errors"
	"github.com/stas137/ms-3/inventory/internal/model"
	inventoryv1 "github.com/stas137/ms-3/shared/pkg/proto/inventory/v1"
)

func TestGet(t *testing.T) {
	t.Parallel()

	type args struct {
		req *inventoryv1.GetPartRequest
	}

	type expected struct {
		res *inventoryv1.GetPartResponse
		err error
	}

	var (
		ctx           = context.Background()
		fakeUUID      = uuid.MustParse(gofakeit.UUID())
		name          = gofakeit.Name()
		description   = gofakeit.Product().Description
		price         = int64(gofakeit.Price(100, 100000))
		partType      = model.PartTypeEngine
		stockQuantity = int64(10)
		createdAt     = time.Now()
	)

	modelPart := model.Part{
		UUID:          fakeUUID,
		Name:          name,
		Description:   description,
		Price:         price,
		PartType:      partType,
		StockQuantity: stockQuantity,
		CreatedAt:     createdAt,
	}

	tests := []struct {
		name      string
		args      args
		setupMock func(repo *mocks.PartService)
		expected  expected
	}{
		{
			name: "успешное получение детали",
			args: args{req: &inventoryv1.GetPartRequest{
				Uuid: fakeUUID.String(),
			}},
			setupMock: func(repo *mocks.PartService) {
				repo.EXPECT().Get(ctx, fakeUUID).Return(modelPart, nil)
			},
			expected: expected{
				res: &inventoryv1.GetPartResponse{
					Part: converter.PartToDTO(modelPart),
				},
				err: nil,
			},
		},
		{
			name: "ошибка сервисного слоя при получении детали",
			args: args{req: &inventoryv1.GetPartRequest{
				Uuid: uuid.Nil.String(),
			}},
			setupMock: func(repo *mocks.PartService) {
				repo.EXPECT().Get(ctx, uuid.Nil).Return(model.Part{}, errs.ErrInvalidUUID)
			},
			expected: expected{
				res: &inventoryv1.GetPartResponse{},
				err: errs.ErrInvalidUUID,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			partService := mocks.NewPartService(t)
			tc.setupMock(partService)

			api := inventory.NewApi(partService)
			res, err := api.GetPart(ctx, tc.args.req)

			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				assert.Equal(t, &inventoryv1.GetPartResponse{}, res)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected.res, res)
			}
		})
	}
}
