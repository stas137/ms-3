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
	"github.com/stas137/ms-3/inventory/internal/model"
	"github.com/stas137/ms-3/inventory/internal/service/input"
	inventoryv1 "github.com/stas137/ms-3/shared/pkg/proto/inventory/v1"
)

func TestList(t *testing.T) {
	t.Parallel()

	type args struct {
		req *inventoryv1.ListPartsRequest
	}

	type expected struct {
		res *inventoryv1.ListPartsResponse
		err error
	}

	var (
		ctx               = context.Background()
		fakeUUID          = uuid.MustParse(gofakeit.UUID())
		name              = gofakeit.Name()
		description       = gofakeit.Product().Description
		price             = int64(gofakeit.Price(100, 100000))
		partType          = model.PartTypeEngine
		partTypeInventory = inventoryv1.PartType_PART_TYPE_ENGINE
		stockQuantity     = int64(10)
		createdAt         = time.Now()
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
			name: "успешное получение деталей",
			args: args{req: &inventoryv1.ListPartsRequest{
				Uuids:    []string{fakeUUID.String()},
				PartType: partTypeInventory,
			}},
			setupMock: func(repo *mocks.PartService) {
				repo.EXPECT().List(ctx, input.PartFilter{
					UUIDs:    []uuid.UUID{fakeUUID},
					PartType: converter.PartTypeToModelPartType(partTypeInventory),
				}).Return([]model.Part{modelPart}, nil)
			},
			expected: expected{
				res: &inventoryv1.ListPartsResponse{
					Parts: []*inventoryv1.Part{
						converter.PartToDTO(modelPart),
					},
				},
				err: nil,
			},
		},
		// {
		// 	name: "ошибка репозитория при получение детали",
		// 	args: args{argUUID: fakeUUID},
		// 	setupMock: func(repo *mocks.PartRepository) {
		// 		repo.EXPECT().Get(ctx, fakeUUID).Return(model.Part{}, errs.ErrPartNotFound)
		// 	},
		// 	expected: expected{
		// 		part: model.Part{},
		// 		err:  errs.ErrPartNotFound,
		// 	},
		// },
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			partService := mocks.NewPartService(t)
			tc.setupMock(partService)

			api := inventory.NewApi(partService)
			res, err := api.ListParts(ctx, tc.args.req)

			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				assert.Equal(t, nil, res)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected.res, res)
			}
		})
	}
}
