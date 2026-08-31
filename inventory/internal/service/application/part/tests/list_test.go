package tests

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/stas137/ms-3/inventory/internal/model"
	"github.com/stas137/ms-3/inventory/internal/service/application/part"
	"github.com/stas137/ms-3/inventory/internal/service/application/part/mocks"
	"github.com/stas137/ms-3/inventory/internal/service/input"
)

func TestList(t *testing.T) {
	t.Parallel()

	type args struct {
		filter input.PartFilter
	}

	type expected struct {
		parts []model.Part
		err   error
	}

	var (
		ctx           = context.Background()
		fakeUUID      = uuid.MustParse(gofakeit.UUID())
		name          = gofakeit.Name()
		description   = gofakeit.Product().Description
		price         = int64(gofakeit.Price(100, 100000))
		partType      = model.PartTypeEngine
		stockQuantity = 10
		reserved      = 5
		createdAt     = time.Now()
	)

	modelPart := model.RestorePart(
		fakeUUID,
		name,
		description,
		partType,
		price,
		stockQuantity,
		reserved,
		model.PartProperties{},
		createdAt,
	)

	tests := []struct {
		name      string
		args      args
		setupMock func(repo *mocks.PartRepository)
		expected  expected
	}{
		{
			name: "успешное получение списка деталей по uuids",
			args: args{filter: input.PartFilter{
				UUIDs:    []uuid.UUID{fakeUUID},
				PartType: partType,
			}},
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().List(
					ctx,
					[]uuid.UUID{fakeUUID},
					partType,
				).
					Return([]model.Part{modelPart}, nil)
			},
			expected: expected{
				parts: []model.Part{modelPart},
				err:   nil,
			},
		}, {
			name: "успешное получение списка деталей по partType",
			args: args{filter: input.PartFilter{
				UUIDs:    []uuid.UUID{},
				PartType: partType,
			}},
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().List(
					ctx,
					[]uuid.UUID{},
					partType,
				).
					Return([]model.Part{modelPart}, nil)
			},
			expected: expected{
				parts: []model.Part{modelPart},
				err:   nil,
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

			partRepository := mocks.NewPartRepository(t)
			compabilityChecker := mocks.NewCompatibilityChecker(t)
			txManager := mocks.NewTxManager(t)
			tc.setupMock(partRepository)

			svc := part.NewService(partRepository, compabilityChecker, txManager)
			res, err := svc.List(ctx, tc.args.filter)

			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				assert.Equal(t, nil, res)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected.parts, res)
			}
		})
	}
}
