package app

import (
	"log/slog"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	orderv1API "github.com/stas137/ms-3/order/internal/api/order/v1"
	inventoryClient "github.com/stas137/ms-3/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/stas137/ms-3/order/internal/client/grpc/payment/v1"
	orderRepository "github.com/stas137/ms-3/order/internal/repository/order"
	orderService "github.com/stas137/ms-3/order/internal/service/order"
	orderv1 "github.com/stas137/ms-3/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/stas137/ms-3/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/stas137/ms-3/shared/pkg/proto/payment/v1"
)

const (
	inventoryServiceAddress = "localhost:50051"
	paymentServiceAddress   = "localhost:50052"
)

func NewServer() *orderv1.Server {
	// Создаём хранилище
	orderRepo := orderRepository.NewRepository()

	inventoryConn, err := grpc.NewClient(inventoryServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             3 * time.Second,
			PermitWithoutStream: true,
		}))
	if err != nil {
		slog.Error("не удалось подключиться к InventoryService", "error", err)
		// return err
		// os.Exit(1)
	}
	defer func() {
		err := inventoryConn.Close()
		if err != nil {
			slog.Error("ошибка закрытия gRPC соединения", "error", err)
		}
	}()

	paymentConn, err := grpc.NewClient(paymentServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             3 * time.Second,
			PermitWithoutStream: true,
		}))
	if err != nil {
		slog.Error("не удалось подключиться к PaymentService", "error", err)
		// return err
		// os.Exit(1)
	}
	defer func() {
		err := paymentConn.Close()
		if err != nil {
			slog.Error("ошибка закрытия gRPC соединения", "error", err)
		}
	}()

	orderInventoryClient := inventoryClient.New(inventoryv1.NewInventoryServiceClient(inventoryConn))
	orderPaymentClient := paymentClient.New(paymentv1.NewPaymentServiceClient(paymentConn))
	orderServ := orderService.NewService(orderRepo, orderInventoryClient, orderPaymentClient)
	orderApi := orderv1API.NewApi(orderServ)

	// h := orderHandler.NewHandler(
	// 	inventoryv1.NewInventoryServiceClient(inventoryConn),
	// 	paymentv1.NewPaymentServiceClient(paymentConn),
	// 	store,
	// )

	// Сгенерировать код ogen из OpenAPI спецификации
	// Команда: task ogen:gen

	// Создать OpenAPI сервер
	// orderServer, err := orderHandler.SetupServer(h)
	// if err != nil {
	// 	slog.Error("ошибка создания сервера OpenAPI", "error", err)
	// 	return err
	// 	// os.Exit(1)
	// }

	// server, err := orderv1.NewServer(apiHandler, orderv1.WithErrorHandler(orderv1API.ErrorHandler))
	server, err := orderv1.NewServer(orderApi, orderv1.WithErrorHandler(orderv1API.ErrorHandler))
	if err != nil {
		slog.Error("ошибка создания сервера OpenAPI", "error", err)
	}

	return server
}

func NewHTTPHandler(
	inventoryServiceClient inventoryv1.InventoryServiceClient,
	paymentServiceClient paymentv1.PaymentServiceClient,
) (http.Handler, error) {
	orderInventoryClient := inventoryClient.New(inventoryServiceClient)
	orderPaymentClient := paymentClient.New(paymentServiceClient)

	// Создаём хранилище
	orderRepo := orderRepository.NewRepository()
	orderServ := orderService.NewService(orderRepo, orderInventoryClient, orderPaymentClient)
	orderApi := orderv1API.NewApi(orderServ)

	return orderv1.NewServer(orderApi, orderv1.WithErrorHandler(orderv1API.ErrorHandler))
}
