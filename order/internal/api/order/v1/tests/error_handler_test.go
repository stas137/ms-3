package tests

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/stretchr/testify/assert"

	"github.com/stas137/ms-3/order/internal/api/order/v1"
	errs "github.com/stas137/ms-3/order/internal/errors"
)

type errorResponse struct {
	Code    int
	Message string
}

func TestErrorHandler(t *testing.T) {
	type args struct {
		err error
	}

	type expected struct {
		code int
		msg  string
		ct   string
	}

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "заказ не найден (errs.ErrOrderNotFound)",
			args: args{
				err: errs.ErrOrderNotFound,
			},
			expected: expected{
				code: 404,
				msg:  errs.ErrOrderNotFound.Error(),
				ct:   "application/json",
			},
		}, {
			name: "деталь не найдена (errs.ErrPartNotFound)",
			args: args{
				err: errs.ErrPartNotFound,
			},
			expected: expected{
				code: 404,
				msg:  errs.ErrPartNotFound.Error(),
				ct:   "application/json",
			},
		}, {
			name: "заказ уже оплачен (errs.ErrOrderAlreadyPaid)",
			args: args{
				err: errs.ErrOrderAlreadyPaid,
			},
			expected: expected{
				code: 409,
				msg:  errs.ErrOrderAlreadyPaid.Error(),
				ct:   "application/json",
			},
		}, {
			name: "заказ уже отменен (errs.ErrOrderCancelled)",
			args: args{
				err: errs.ErrOrderCancelled,
			},
			expected: expected{
				code: 409,
				msg:  errs.ErrOrderCancelled.Error(),
				ct:   "application/json",
			},
		}, {
			name: "заказ с невалидным uuid (errs.ErrInvalidUUID)",
			args: args{
				err: errs.ErrInvalidUUID,
			},
			expected: expected{
				code: 400,
				msg:  errs.ErrInvalidUUID.Error(),
				ct:   "application/json",
			},
		}, {
			name: "заказ с невалидным uuid (errs.ErrInvalidPaymentMethod)",
			args: args{
				err: errs.ErrInvalidPaymentMethod,
			},
			expected: expected{
				code: 400,
				msg:  errs.ErrInvalidPaymentMethod.Error(),
				ct:   "application/json",
			},
		}, {
			name: "заказ (внутренняя ошибка)",
			args: args{
				err: errors.New("внутренняя ошибка"),
			},
			expected: expected{
				code: 500,
				msg:  errors.New("внутренняя ошибка").Error(),
				ct:   "application/json",
			},
		}, {
			name: "заказ (ошибка DecodeRequestError)",
			args: args{
				err: &ogenerrors.DecodeRequestError{},
			},
			expected: expected{
				code: 400,
				msg:  (&ogenerrors.DecodeRequestError{}).Error(),
				ct:   "application/json",
			},
		}, {
			name: "заказ (ошибка DecodeParamsError)",
			args: args{
				err: &ogenerrors.DecodeParamsError{},
			},
			expected: expected{
				code: 400,
				msg:  (&ogenerrors.DecodeParamsError{}).Error(),
				ct:   "application/json",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()

			order.ErrorHandler(context.Background(), rec, nil, tc.args.err)

			// статус код
			if got := rec.Code; got != tc.expected.code {
				t.Errorf("статус = %d, ожидался = %d", got, tc.expected.code)
			} else {
				assert.Equal(t, tc.expected.code, rec.Code)
			}

			// Content-Type
			if ct := rec.Header().Get("Content-Type"); ct != tc.expected.ct {
				t.Errorf(
					"Content-Type = %q, ожидался application/json",
					ct,
				)
			}

			// Тело ответа
			var resp errorResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
				t.Fatalf("не удалось декодировать JSON: %v", err)
			}
			if resp.Code != tc.expected.code {
				t.Errorf("json code = %d, ожидался %d", resp.Code, tc.expected.code)
			}
			if resp.Message != tc.expected.msg {
				t.Errorf("json message = %q, ожидалось %q", resp.Message, tc.expected.msg)
			}
		})
	}
}
