package main

import (
	"context"
	"log/slog"

	"github.com/joho/godotenv"

	"github.com/stas137/ms-3/inventory/internal/app"
	"github.com/stas137/ms-3/inventory/internal/config"
)

// const (
// 	grpcAddress               = ":50051"
// 	grpcMaxConnectionIdle     = 15 * time.Minute // Закрыть idle-соединение (нет активных RPC)
// 	grpcMaxConnectionAge      = 30 * time.Minute // Принудительная ротация для балансировки
// 	grpcMaxConnectionAgeGrace = 5 * time.Second  // Время на завершение активных RPC
// 	grpcKeepaliveTime         = 5 * time.Minute  // Интервал ping'ов для обнаружения мертвых соединений
// 	grpcKeepaliveTimeout      = 1 * time.Second  // Таймаут ожидания pong
// 	grpcMinPingInterval       = 5 * time.Minute  // Минимальный интервал ping'ов от клиента (защита от DoS)
// )

func main() {
	err := godotenv.Load("../inventory.env")
	if err != nil {
		slog.Error("ошибка при загрузке inventory.env", "error", err)
	}

	configPath := config.ResolveConfigPath()
	config.MustLoad(configPath)

	a := app.New(context.Background())

	if err := a.Run(); err != nil {
		slog.Error("ошибка при работе приложения", "error", err)
		// os.Exit(1)
	}
}

// func run() error {
// 	ctx, cancel := context.WithTimeout(context.Background(), grpcMaxConnectionAge)
// 	defer cancel()

// 	_ = godotenv.Load("../inventory.env")
// 	configPath := config.ResolveConfigPath()

// 	cfg, err := config.Load(configPath)
// 	if err != nil {
// 		return err
// 	}

// 	// postgresDB := os.Getenv("POSTGRES_DB")

// 	slog.Info("конфигурация загружена",
// 		"config_path", configPath,
// 		"grpc_address", cfg.GRPC.Address(),
// 		"pg_host", cfg.PG.Host,
// 		"DSN", cfg.PG.DSN(),
// 		// "postgresDB", postgresDB,
// 	)

// 	lis, err := (*net.ListenConfig).Listen(&net.ListenConfig{}, ctx, "tcp", string(cfg.GRPC.Address()))
// 	// lis, err := net.Listen("tcp", grpcAddress)
// 	if err != nil {
// 		return err
// 	}

// 	// Настроить gRPC сервер с параметрами keepalive
// 	// Подумайте, какие параметры стоит задать для production-ready сервера
// 	// См. examples/week_1/GRPC_CONNECTIONS.md
// 	grpcServer := grpc.NewServer(
// 		grpc.KeepaliveParams(keepalive.ServerParameters{
// 			MaxConnectionIdle:     grpcMaxConnectionIdle,
// 			MaxConnectionAge:      grpcMaxConnectionAge,
// 			MaxConnectionAgeGrace: grpcMaxConnectionAgeGrace,
// 			Time:                  grpcKeepaliveTime,
// 			Timeout:               grpcKeepaliveTimeout,
// 		}),
// 		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
// 			MinTime:             grpcMinPingInterval,
// 			PermitWithoutStream: true, // Разрешить "теплые" соединения без активных RPC
// 		}),
// 	)

// 	// inventoryv1.RegisterInventoryServiceServer(grpcServer, inventoryService.NewServer())
// 	// inventoryv1.RegisterInventoryServiceServer(grpcServer, app.NewServer())
// 	// dbURI := os.Getenv("DB_URI")

// 	pool, err := pgxpool.New(ctx, cfg.PG.DSN())
// 	if err != nil {
// 		return err
// 	}
// 	defer pool.Close()

// 	err = pool.Ping(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	slog.Info("подключение к Postgres установлено")

// 	app.RegisterServices(grpcServer, pool)

// 	// Включаем reflection для postman/grpcurl
// 	reflection.Register(grpcServer)

// 	slog.Info("запуск InventoryService", "адрес", grpcAddress)

// 	// Реализовать graceful shutdown
// 	// При получении сигнала SIGINT/SIGTERM сервер должен:
// 	// 1. Перестать принимать новые соединения
// 	// 2. Дождаться завершения текущих запросов
// 	// 3. Корректно завершить работу
// 	// Подсказка: используйте signal.NotifyContext и grpcServer.GracefulStop()

// 	go func() {
// 		slog.Info("inventory gRPC сервер запущен", "адрес", grpcAddress)
// 		if serveErr := grpcServer.Serve(lis); serveErr != nil {
// 			slog.Error("ошибка запуска сервера", "error", serveErr)
// 		}
// 	}()

// 	// Graceful shutdown
// 	// quit := make(chan os.Signal, 1)
// 	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
// 	// <-quit

// 	// ctx, cancel := context.WithCancel(context.Background())
// 	// defer cancel()

// 	ctx, cancel = signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
// 	defer cancel()

// 	<-ctx.Done()
// 	slog.Info("остановка gRPC сервера")
// 	grpcServer.GracefulStop()
// 	slog.Info("сервер остановлен")

// 	return nil
// }
