package order

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	errs "github.com/stas137/ms-3/order/internal/errors"
	"github.com/stas137/ms-3/order/internal/model"
	repoConverter "github.com/stas137/ms-3/order/internal/repository/converter"
	"github.com/stas137/ms-3/order/internal/repository/record"
)

type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

// хранилище заказов (in-memory) потокобезовасное
type repository struct {
	orderPool *pgxpool.Pool
	getter    *trmpgx.CtxGetter
	txManager TxManager
}

// NewRepository создаёт новое пустое хранилище заказов
func NewRepository(orderPool *pgxpool.Pool, txManager TxManager) *repository {
	return &repository{
		orderPool: orderPool,
		getter:    trmpgx.DefaultCtxGetter,
		txManager: txManager,
	}
}

func (r *repository) getOrder(ctx context.Context, orderUUID uuid.UUID) (model.Order, error) {
	var orderDTO record.Order
	var orderItemsDTO []record.OrderItem

	queryOrder := `SELECT uuid, status, transaction_uuid, payment_method, created_at, updated_at FROM orders WHERE uuid = $1`
	queryOrderItems := `SELECT order_uuid, part_uuid, part_type, price FROM order_items WHERE order_uuid = $1`

	row := r.getter.DefaultTrOrDB(ctx, r.orderPool).QueryRow(ctx, queryOrder, orderUUID)

	err := row.Scan(
		&orderDTO.UUID,
		&orderDTO.Status,
		&orderDTO.TransactionUUID,
		&orderDTO.PaymentMethod,
		&orderDTO.CreatedAt,
		&orderDTO.UpdatedAt,
	)
	if err != nil {
		return model.Order{}, errs.ErrOrderNotFound
	}

	rows, err := r.getter.DefaultTrOrDB(ctx, r.orderPool).Query(ctx, queryOrderItems, orderUUID)
	if err != nil {
		return model.Order{}, errs.ErrPartNotFound
	}
	defer rows.Close()

	orderItemsDTO, err = pgx.CollectRows(rows, pgx.RowToStructByPos[record.OrderItem])
	if err != nil {
		return model.Order{}, fmt.Errorf("считать строки: %w", err)
	}

	return repoConverter.RepoOrderToModelOrder(orderDTO, orderItemsDTO), nil
}

func (r *repository) createOrder(ctx context.Context, order model.Order) error {
	rec := repoConverter.ModelOrderToRepoOrder(order)

	const query = `
		INSERT INTO orders (uuid, status, created_at)
		VALUES ($1, $2, $3)
	`

	_, err := r.getter.DefaultTrOrDB(ctx, r.orderPool).Exec(ctx, query,
		rec.UUID,
		rec.Status,
		time.Now())
	if err != nil {
		return fmt.Errorf("создать заказ: %w", err)
	}

	return nil
}

func (r *repository) createOrderItems(ctx context.Context, order model.Order) error {
	if len(order.Items) == 0 {
		return nil
	}

	items := repoConverter.ModelOrderItemsToRepoOrderItems(order)

	query := squirrel.Insert("order_items").Columns("order_uuid", "part_uuid", "part_type", "price").
		PlaceholderFormat(squirrel.Dollar)

	for _, item := range items {
		query = query.Values(item.OrderUUID, item.PartUUID, item.PartType, item.Price)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	_, err = r.getter.DefaultTrOrDB(ctx, r.orderPool).Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("insert order items: %w", err)
	}

	return nil
}

func (r *repository) updateOrder(ctx context.Context, order model.Order) error {
	const queryOrder = `
		UPDATE orders SET 
			payment_method =	COALESCE($1, payment_method), 
			status = 			COALESCE($2, status), 
			transaction_uuid = 	COALESCE($3, transaction_uuid), 
			updated_at = 		$4 
		WHERE uuid = $5;`

	res, err := r.getter.DefaultTrOrDB(ctx, r.orderPool).Exec(ctx, queryOrder,
		order.PaymentMethod,
		order.Status,
		order.TransactionUUID,
		time.Now(),
		order.UUID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errs.ErrOrderNotFound
	}

	return nil
}

func (r *repository) updateOrderItems(ctx context.Context, order model.Order) error {
	const queryOrderItems = `
		UPDATE order_items SET
			part_type = 		COALESCE($1, part_type),
			price = 			COALESCE($2, price)
			WHERE order_uuid = $3 AND part_uuid = $4`

	for _, item := range order.Items {
		res, err := r.getter.DefaultTrOrDB(ctx, r.orderPool).Exec(ctx, queryOrderItems,
			item.PartType,
			item.Price,
			order.UUID,
			item.PartUUID)
		if err != nil {
			return err
		}
		if res.RowsAffected() == 0 {
			return errs.ErrPartNotFound
		}
	}
	return nil
}

func (r *repository) Get(ctx context.Context, orderUUID uuid.UUID) (model.Order, error) {
	return r.getOrder(ctx, orderUUID)
}

func (r *repository) Create(ctx context.Context, order model.Order) error {
	return r.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := r.createOrder(txCtx, order); err != nil {
			return err
		}
		return r.createOrderItems(txCtx, order)
	})
}

func (r *repository) Update(ctx context.Context, order model.Order) error {
	return r.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := r.updateOrder(txCtx, order); err != nil {
			return err
		}
		return r.updateOrderItems(txCtx, order)
	})
}
