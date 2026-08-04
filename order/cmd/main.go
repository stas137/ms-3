package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"github.com/stas137/ms-3/order/pkg/app"
)

const (
	inventoryServiceAddress = "localhost:50051"
	paymentServiceAddress   = "localhost:50052"
	httpPort                = "8080"
	readHeaderTimeout       = 5 * time.Second
	readTimeout             = 15 * time.Second
	writeTimeout            = 15 * time.Second
	idleTimeout             = 60 * time.Second
	maxHeaderBytes          = 1 << 20
	shutdownTimeout         = 10 * time.Second
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	// Настроить gRPC клиент с параметрами keepalive
	// Подумайте, какие параметры стоит задать для gRPC клиента
	// См. examples/week_1/GRPC_CONNECTIONS.md

	// Создать gRPC соединение с InventoryService
	inventoryConn, err := grpc.NewClient(inventoryServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             3 * time.Second,
			PermitWithoutStream: true,
		}))
	if err != nil {
		slog.Error("не удалось подключиться к InventoryService", "error", err)
		return err
		// os.Exit(1)
	}
	defer func() {
		err := inventoryConn.Close()
		if err != nil {
			slog.Error("ошибка закрытия gRPC соединения", "error", err)
		}
	}()

	// Создать gRPC клиент PaymentService
	paymentConn, err := grpc.NewClient(paymentServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             3 * time.Second,
			PermitWithoutStream: true,
		}))
	if err != nil {
		slog.Error("не удалось подключиться к PaymentService", "error", err)
		return err
		// os.Exit(1)
	}
	defer func() {
		err := paymentConn.Close()
		if err != nil {
			slog.Error("ошибка закрытия gRPC соединения", "error", err)
		}
	}()

	// // Создаём хранилище и обработчик
	// store := orderHandler.NewOrderStore()
	// h := orderHandler.NewHandler(
	// 	inventoryv1.NewInventoryServiceClient(inventoryConn),
	// 	paymentv1.NewPaymentServiceClient(paymentConn),
	// 	store,
	// )

	// // Сгенерировать код ogen из OpenAPI спецификации
	// // Команда: task ogen:gen

	// // Создать OpenAPI сервер
	// orderServer, err := orderHandler.SetupServer(h)
	// if err != nil {
	// 	slog.Error("ошибка создания сервера OpenAPI", "error", err)
	// 	return err
	// 	// os.Exit(1)
	// }

	// Настроить HTTP сервер с таймаутами
	// Создайте &http.Server{...} с явными таймаутами вместо http.ListenAndServe(...)
	// Минимальный набор: ReadHeaderTimeout (защита от Slowloris), ReadTimeout, WriteTimeout, IdleTimeout
	// Без ReadHeaderTimeout сервер уязвим к атаке Slowloris (медленная отправка заголовков)
	// См. examples/week_1/HTTP_SERVER.md

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           app.NewServer(),
		ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атаки
		ReadTimeout:       readTimeout,       // Лимит на чтение всего запроса
		WriteTimeout:      writeTimeout,      // Лимит на запись ответа
		IdleTimeout:       idleTimeout,       // Таймаут keep-alive соединения
		MaxHeaderBytes:    maxHeaderBytes,    // 1 MB // Максимальный размер заголовка
	}

	// Реализовать graceful shutdown для HTTP сервера
	// При получении сигнала SIGINT/SIGTERM сервер должен:
	// 1. Перестать принимать новые соединения
	// 2. Дождаться завершения текущих запросов (с таймаутом)
	// 3. Закрыть gRPC соединения
	// 4. Корректно завершить работу
	// Подсказка: используйте signal.NotifyContext и httpServer.Shutdown(ctx)

	slog.Info("запуск OrderService", "port", 8080)

	// Запускаем сервер в отдельной горутине
	go func() {
		slog.Info("HTTP-сервер запущен на порту", "port", httpPort)
		listenErr := server.ListenAndServe()
		if listenErr != nil && !errors.Is(listenErr, http.ErrServerClosed) {
			slog.Error("ошибка запуска сервера", "error", listenErr)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("ошибка при остановке сервера", "error", err)
	}

	slog.Info("сервер оставновлен")

	return nil

	// err = http.ListenAndServe(":8080", orderServer)
	// if err != nil {
	// 	slog.Error("ошибка запуска сервера", "error", err)
	// 	os.Exit(1)
	// }
}
