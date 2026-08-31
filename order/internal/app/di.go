package app

import (
	"context"
	"log/slog"
	"os"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	orderv1API "github.com/stas137/ms-3/order/internal/api/order/v1"
	inventoryClient "github.com/stas137/ms-3/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/stas137/ms-3/order/internal/client/grpc/payment/v1"
	"github.com/stas137/ms-3/order/internal/config"
	orderRepository "github.com/stas137/ms-3/order/internal/repository/order"
	orderService "github.com/stas137/ms-3/order/internal/service/order"
	"github.com/stas137/ms-3/platform/pkg/closer"
	orderv1 "github.com/stas137/ms-3/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/stas137/ms-3/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/stas137/ms-3/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	pgPool    *pgxpool.Pool
	txManager orderService.TxManager

	inventoryConn *grpc.ClientConn
	paymentConn   *grpc.ClientConn

	inventoryClient orderService.InventoryClient
	paymentClient   orderService.PaymentClient

	orderRepo orderService.OrderRepository
	orderServ orderv1API.OrderService
	orderAPI  orderv1.Handler

	orderServer *orderv1.Server
}

func (d *diContainer) PGPool(ctx context.Context) *pgxpool.Pool {
	if d.pgPool == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().PG.DSN())
		if err != nil {
			slog.Error("не удалось подключиться к PostgreSQL", "error", err)
			os.Exit(1)
		}

		err = pool.Ping(ctx)
		if err != nil {
			slog.Error("не удалось выполнить ping PostgreSQL", "error", err)
			os.Exit(1)
		}

		closer.Add("PostgreSQL pool", func(_ context.Context) error {
			pool.Close()
			return nil
		})
		d.pgPool = pool
	}
	return d.pgPool
}

func (d *diContainer) TxManager(ctx context.Context) orderService.TxManager {
	if d.txManager != nil {
		m, err := manager.New(trmpgx.NewDefaultFactory(d.PGPool(ctx)))
		if err != nil {
			slog.Error("не удалось создать transaction manager", "error", err)
			os.Exit(1)
		}

		d.txManager = m
	}
	return d.txManager
}

func (d *diContainer) OrderRepo(ctx context.Context) orderService.OrderRepository {
	if d.orderRepo == nil {
		orderRepo := orderRepository.NewRepository(d.PGPool(ctx), d.TxManager(ctx))
		d.orderRepo = orderRepo
	}
	return d.orderRepo
}

func (d *diContainer) InventoryConn(_ context.Context) *grpc.ClientConn {
	if d.inventoryConn == nil {
		invConn, err := grpc.NewClient(config.AppConfig().InventoryClient.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                10 * time.Second,
				Timeout:             3 * time.Second,
				PermitWithoutStream: true,
			}))
		if err != nil {
			slog.Error("не удалось подключиться к InventoryService", "error", err)
			// return err
			os.Exit(1)
		}
		d.inventoryConn = invConn
	}
	return d.inventoryConn
}

func (d *diContainer) PaymentConn(_ context.Context) *grpc.ClientConn {
	if d.paymentConn == nil {
		paymentConn, err := grpc.NewClient(config.AppConfig().PaymentClient.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                10 * time.Second,
				Timeout:             3 * time.Second,
				PermitWithoutStream: true,
			}))
		if err != nil {
			slog.Error("не удалось подключиться к PaymentService", "error", err)
			// return err
			os.Exit(1)
		}
		d.paymentConn = paymentConn
	}
	return d.paymentConn
}

func (d *diContainer) InventoryClient(ctx context.Context) orderService.InventoryClient {
	if d.inventoryClient == nil {
		invClient := inventoryClient.New(inventoryv1.NewInventoryServiceClient(d.InventoryConn(ctx)))
		d.inventoryClient = invClient
	}
	return d.inventoryClient
}

func (d *diContainer) PaymentClient(ctx context.Context) orderService.PaymentClient {
	if d.paymentClient == nil {
		paymentClient := paymentClient.New(paymentv1.NewPaymentServiceClient(d.PaymentConn(ctx)))
		d.paymentClient = paymentClient
	}
	return d.paymentClient
}

func (d *diContainer) OrderServ(ctx context.Context) orderv1API.OrderService {
	if d.orderServ == nil {
		orderServ := orderService.NewService(d.OrderRepo(ctx), d.InventoryClient(ctx), d.PaymentClient(ctx), d.TxManager(ctx))
		d.orderServ = orderServ
	}
	return d.orderServ
}

func (d *diContainer) OrderAPI(ctx context.Context) orderv1.Handler {
	if d.orderAPI == nil {
		orderAPI := orderv1API.NewApi(d.OrderServ(ctx))
		d.orderAPI = orderAPI
	}
	return d.orderAPI
}

func (d *diContainer) OrderServer(ctx context.Context) *orderv1.Server {
	if d.orderServer == nil {
		orderServer, err := orderv1.NewServer(d.OrderAPI(ctx), orderv1.WithErrorHandler(orderv1API.ErrorHandler))
		if err != nil {
			slog.Error("ошибка создания сервера OpenAPI", "error", err)
			os.Exit(1)
		}
		d.orderServer = orderServer
	}
	return d.orderServer
}
