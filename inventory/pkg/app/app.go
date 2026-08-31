package app

import (
	"log/slog"
	"os"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	partv1API "github.com/stas137/ms-3/inventory/internal/api/inventory/v1"
	"github.com/stas137/ms-3/inventory/internal/interceptor"
	partRepository "github.com/stas137/ms-3/inventory/internal/repository/part"
	partService "github.com/stas137/ms-3/inventory/internal/service/application/part"
	"github.com/stas137/ms-3/inventory/internal/service/domain"
	inventoryv1 "github.com/stas137/ms-3/shared/pkg/proto/inventory/v1"
)

func RegisterServices(grpcServer *grpc.Server, inventoryPool *pgxpool.Pool) {
	txManager, err := manager.New(trmpgx.NewDefaultFactory(inventoryPool))
	if err != nil {
		slog.Error("не удалось создать transaction manager", "error", err)
		os.Exit(1)
	}
	partRepo := partRepository.NewRepository(inventoryPool)
	partServ := partService.NewService(partRepo, domain.NewCompatibilityChecker(), txManager)
	partApi := partv1API.NewApi(partServ)

	inventoryv1.RegisterInventoryServiceServer(grpcServer, partApi)
}

func Interceptors() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.ErrorInterceptor),
	}
}
