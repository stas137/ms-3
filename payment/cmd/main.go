package main

import (
	"context"
	"log/slog"

	"github.com/joho/godotenv"

	"github.com/stas137/ms-3/payment/internal/app"
	"github.com/stas137/ms-3/payment/internal/config"
)

func main() {
	err := godotenv.Load("../payment.env")
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

	// if err := run(); err != nil {
	// 	slog.Error("не удалось создать payment listener", "error", err)
	// 	os.Exit(1)
	// }

	// err = grpcServer.Serve(lis)
	// if err != nil {
	// 	slog.Error("ошибка запуска сервера", "error", err)
	// 	os.Exit(1)
	// }
}

// func run() error {
// 	ctx, cancel := context.WithTimeout(context.Background(), grpcMaxConnectionAge)
// 	defer cancel()

// 	lis, err := (*net.ListenConfig).Listen(&net.ListenConfig{}, ctx, "tcp", string(grpcAddress))
// 	if err != nil {
// 		return err
// 	}

// 	// Настроить gRPC сервер с параметрами keepalive
// 	// Подумайте, какие параметры стоит задать для production-ready сервера
// 	// См. examples/week_1/GRPC_CONNECTIONS.md
// 	grpcServer := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
// 		MaxConnectionIdle:     grpcMaxConnectionIdle,
// 		MaxConnectionAge:      grpcMaxConnectionAge,
// 		MaxConnectionAgeGrace: grpcMaxConnectionAgeGrace,
// 		Time:                  grpcKeepaliveTime,
// 		Timeout:               grpcKeepaliveTimeout,
// 	}), grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
// 		MinTime:             grpcMinPingInterval,
// 		PermitWithoutStream: true,
// 	}))

// 	// paymentv1.RegisterPaymentServiceServer(grpcServer, paymentService.NewServer())
// 	app.RegisterServices(grpcServer)

// 	// Включаем reflection для postman/grpcurl
// 	reflection.Register(grpcServer)

// 	slog.Info("запуск PaymentService", "адрес", grpcAddress)

// 	// Реализовать graceful shutdown
// 	// При получении сигнала SIGINT/SIGTERM сервер должен:
// 	// 1. Перестать принимать новые соединения
// 	// 2. Дождаться завершения текущих запросов
// 	// 3. Корректно завершить работу
// 	// Подсказка: используйте signal.NotifyContext и grpcServer.GracefulStop()

// 	go func() {
// 		slog.Info("payment gRPC сервер запущен", "адрес", grpcAddress)
// 		if serveError := grpcServer.Serve(lis); serveError != nil {
// 			slog.Error("ошибка запуска сервера", "error", serveError)
// 		}
// 	}()

// 	// Graceful shutdown
// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
// 	<-quit

// 	slog.Info("остановка gRPC сервера")
// 	grpcServer.GracefulStop()
// 	slog.Info("сервер остановлен")

// 	return nil
// }
