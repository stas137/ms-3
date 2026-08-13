package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	partv1API "github.com/stas137/ms-3/inventory/internal/api/inventory/v1"
	"github.com/stas137/ms-3/inventory/internal/interceptor"
	partRepository "github.com/stas137/ms-3/inventory/internal/repository/part"
	partService "github.com/stas137/ms-3/inventory/internal/service/part"
	inventoryv1 "github.com/stas137/ms-3/shared/pkg/proto/inventory/v1"
)

func RegisterServices(grpcServer *grpc.Server, inventoryPool *pgxpool.Pool) {
	partRepo := partRepository.NewRepository(inventoryPool)
	partServ := partService.NewService(partRepo)
	partApi := partv1API.NewApi(partServ)

	inventoryv1.RegisterInventoryServiceServer(grpcServer, partApi)
}

func Interceptors() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.ErrorInterceptor),
	}
}
