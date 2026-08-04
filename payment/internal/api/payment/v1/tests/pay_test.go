package tests

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/stas137/ms-3/payment/internal/api/payment/v1"
	"github.com/stas137/ms-3/payment/internal/api/payment/v1/mocks"
	"github.com/stas137/ms-3/payment/internal/model"
	"github.com/stas137/ms-3/payment/internal/service/input"
	paymentv1 "github.com/stas137/ms-3/shared/pkg/proto/payment/v1"
)

func TestPay(t *testing.T) {
	t.Parallel()

	type args struct {
		payReq *paymentv1.PayOrderRequest
	}

	type expected struct {
		payRes *paymentv1.PayOrderResponse
		err    error
	}

	var (
		ctx      = context.Background()
		fakeUUID = uuid.MustParse(gofakeit.UUID())
	)

	payOrderRes := &paymentv1.PayOrderResponse{
		TransactionUuid: fakeUUID.String(),
	}

	tests := []struct {
		name      string
		args      args
		setupMock func(repo *mocks.PaymentService)
		expected  expected
	}{
		{
			name: "упешная оплата заказа",
			args: args{payReq: &paymentv1.PayOrderRequest{
				OrderUuid:     fakeUUID.String(),
				PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
			}},
			setupMock: func(repo *mocks.PaymentService) {
				repo.EXPECT().Pay(ctx, input.PayOrderInput{
					OrderUUID:     fakeUUID,
					PaymentMethod: model.PaymentMethodCreditCard,
				}).Return(fakeUUID, nil)
			},
			expected: expected{
				payRes: payOrderRes,
				err:    nil,
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

			paymentService := mocks.NewPaymentService(t)
			tc.setupMock(paymentService)

			api := payment.NewApi(paymentService)
			res, err := api.PayOrder(ctx, tc.args.payReq)

			if tc.expected.err != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expected.err)
				assert.Equal(t, nil, res)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected.payRes, res)
			}
		})
	}
}
