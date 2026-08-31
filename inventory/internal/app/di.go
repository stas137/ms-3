package app

import (
	"context"
	"log/slog"
	"os"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"

	partV1API "github.com/stas137/ms-3/inventory/internal/api/inventory/v1"
	"github.com/stas137/ms-3/inventory/internal/config"
	partRepository "github.com/stas137/ms-3/inventory/internal/repository/part"
	partService "github.com/stas137/ms-3/inventory/internal/service/application/part"
	"github.com/stas137/ms-3/inventory/internal/service/domain"
	"github.com/stas137/ms-3/platform/pkg/closer"
	inventoryv1 "github.com/stas137/ms-3/shared/pkg/proto/inventory/v1"
)

// diContainer — контейнер зависимостей с ленивой инициализацией
// Каждый геттер проверяет nil, создаёт объект при первом вызове и кэширует
type diContainer struct {
	// Инфраструктура
	pgPool    *pgxpool.Pool
	txManager partService.TxManager

	// Доменный сервис
	checker partService.CompatibilityChecker

	// Репозитории
	partRepo partService.PartRepository

	// Application-сервис
	partServ partV1API.PartService

	// API-обработчик
	partV1Handler inventoryv1.InventoryServiceServer
}

// При первом вызове создает пул, проверяет соединение и регистрирует closer.
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

func (d *diContainer) PartRepository(ctx context.Context) partService.PartRepository {
	if d.partRepo == nil {
		partRepo := partRepository.NewRepository(d.PGPool(ctx))
		d.partRepo = partRepo
	}
	return d.partRepo
}

func (d *diContainer) PartService(ctx context.Context) partV1API.PartService {
	if d.partServ == nil {
		partServ := partService.NewService(d.PartRepository(ctx), d.CompatibilityChecker(), d.TxManager(ctx))
		d.partServ = partServ
	}
	return d.partServ
}

func (d *diContainer) PartV1API(ctx context.Context) inventoryv1.InventoryServiceServer {
	if d.partV1Handler == nil {
		partV1Handler := partV1API.NewApi(d.PartService(ctx))
		d.partV1Handler = partV1Handler
	}
	return d.partV1Handler
}

// TxManager возвращает менеджер транзакций
func (d *diContainer) TxManager(ctx context.Context) partService.TxManager {
	if d.txManager == nil {
		m, err := manager.New(trmpgx.NewDefaultFactory(d.PGPool(ctx)))
		if err != nil {
			slog.Error("не удалось создать transaction manager", "error", err)
			os.Exit(1)
		}

		d.txManager = m
	}

	return d.txManager
}

// // ComponentRepository возвращает репозиторий комплектующих
// func (d *diContainer) ComponentRepository(ctx context.Context) pcBuilder.ComponentRepository {
// 	if d.componentRepo == nil {
// 		d.componentRepo = componentRepo.NewRepository(d.PGPool(ctx))
// 	}

// 	return d.componentRepo
// }

// // BuildRepository возвращает репозиторий сборок
// func (d *diContainer) BuildRepository(ctx context.Context) pcBuilder.BuildRepository {
// 	if d.buildRepo == nil {
// 		d.buildRepo = buildRepo.NewRepository(d.PGPool(ctx))
// 	}

// 	return d.buildRepo
// }

// CompatibilityChecker возвращает доменный сервис проверки совместимости
func (d *diContainer) CompatibilityChecker() partService.CompatibilityChecker {
	if d.checker == nil {
		d.checker = domain.NewCompatibilityChecker()
	}

	return d.checker
}

// // PCBuilderService возвращает application-сервис сборки ПК
// func (d *diContainer) PCBuilderService(ctx context.Context) pcBuilderAPI.PCBuilderService {
// 	if d.pcBuilderSvc == nil {
// 		d.pcBuilderSvc = pcBuilder.NewService(
// 			d.TxManager(ctx),
// 			d.ComponentRepository(ctx),
// 			d.BuildRepository(ctx),
// 			d.CompatibilityChecker(),
// 		)
// 	}

// 	return d.pcBuilderSvc
// }

// // PCBuilderHandler возвращает API-обработчик сборки ПК
// func (d *diContainer) PCBuilderHandler(ctx context.Context) *pcBuilderAPI.Handler {
// 	if d.handler == nil {
// 		d.handler = pcBuilderAPI.NewHandler(d.PCBuilderService(ctx))
// 	}

// 	return d.handler
// }
