package order

import (
	"context"
	"sync"

	"github.com/google/uuid"

	errs "github.com/stas137/ms-3/order/internal/errors"
	"github.com/stas137/ms-3/order/internal/model"
	repoConverter "github.com/stas137/ms-3/order/internal/repository/converter"
	"github.com/stas137/ms-3/order/internal/repository/record"
)

// хранилище заказов (in-memory) потокобезовасное
type repository struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]record.Order
}

// NewRepository создаёт новое пустое хранилище заказов
func NewRepository() *repository {
	return &repository{
		orders: make(map[uuid.UUID]record.Order),
	}
}

func (r *repository) Create(_ context.Context, order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[order.UUID] = repoConverter.ModelOrderToRepoOrder(order)

	return nil // handle error and return error if it necessary
}

func (r *repository) Get(_ context.Context, orderUUID uuid.UUID) (model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[orderUUID]
	if !ok {
		return model.Order{}, errs.ErrOrderNotFound
	}
	if order.DeletedAt != nil {
		return model.Order{}, errs.ErrOrderNotFound
	}

	return repoConverter.RepoOrderToModelOrder(order), nil
}

func (r *repository) Update(_ context.Context, order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	orderExist, ok := r.orders[order.UUID]
	if !ok {
		return errs.ErrOrderNotFound
	}
	if orderExist.DeletedAt != nil {
		return errs.ErrOrderNotFound
	}

	repoOrder := repoConverter.ModelOrderToRepoOrder(order)
	r.orders[order.UUID] = repoOrder

	return nil
}
