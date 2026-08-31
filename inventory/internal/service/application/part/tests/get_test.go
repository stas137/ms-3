package tests

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/stas137/ms-3/inventory/internal/errors"
	"github.com/stas137/ms-3/inventory/internal/model"
	"github.com/stas137/ms-3/inventory/internal/service/application/part"
	"github.com/stas137/ms-3/inventory/internal/service/application/part/mocks"
)

func TestGet(t *testing.T) {
	t.Parallel()

	type args struct {
		argUUID uuid.UUID
	}

	type expected struct {
		part model.Part
		err  error
	}

	var (
		ctx           = context.Background()
		fakeUUID      = uuid.MustParse(gofakeit.UUID())
		name          = gofakeit.Name()
		description   = gofakeit.Product().Description
		price         = int64(gofakeit.Price(100, 100000))
		partType      = model.PartTypeEngine
		stockQuantity = int(10)
		reserved      = int(5)
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
			name: "упешное получение детали",
			args: args{argUUID: fakeUUID},
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().Get(ctx, fakeUUID).Return(modelPart, nil)
			},
			expected: expected{
				part: modelPart,
				err:  nil,
			},
		}, {
			name: "ошибка репозитория при получение детали",
			args: args{argUUID: fakeUUID},
			setupMock: func(repo *mocks.PartRepository) {
				repo.EXPECT().Get(ctx, fakeUUID).Return(model.Part{}, errs.ErrPartNotFound)
			},
			expected: expected{
				part: model.Part{},
				err:  errs.ErrPartNotFound,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			partRepository := mocks.NewPartRepository(t)
			compabilityChecker := mocks.NewCompatibilityChecker(t)
			txManager := mocks.NewTxManager(t)

			tc.setupMock(partRepository)

			svc := part.NewService(partRepository, compabilityChecker, txManager)
			res, err := svc.Get(ctx, tc.args.argUUID)

			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				assert.Equal(t, model.Part{}, res)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected.part, res)
			}
		})
	}
}

// func TestGetSuccess(t *testing.T) {
// s.partRepository.EXPECT().Get(ctx, fakeUUID).Return(modelPart, nil)

// res, err := s.service.Get(ctx, fakeUUID)
// s.NoError(err)
// s.Equal(modelPart, res)

// }

// func TestGetError(t *testing.T) {
// 	var (
// 		repoErr  = gofakeit.Error() // errs.ErrNotFound
// 		fakeUUID = uuid.MustParse(gofakeit.UUID())
// 	)

// 	s.partRepository.EXPECT().Get(s.ctx, fakeUUID).Return(model.Part{}, repoErr)

// 	res, err := s.service.Get(s.ctx, fakeUUID)
// 	s.Error(err)
// 	s.ErrorIs(err, repoErr)
// 	s.Equal(model.Part{}, res)
// }
