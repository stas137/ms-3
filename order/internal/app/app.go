package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/stas137/ms-3/order/internal/config"
	"github.com/stas137/ms-3/platform/pkg/closer"
	"github.com/stas137/ms-3/platform/pkg/logger"
)

const (
	// inventoryServiceAddress = "localhost:50051"
	// paymentServiceAddress   = "localhost:50052"
	// httpPort                = "8080"
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 15 * time.Second
	writeTimeout      = 15 * time.Second
	idleTimeout       = 60 * time.Second
	maxHeaderBytes    = 1 << 20
	shutdownTimeout   = 10 * time.Second
)

// App - корневая структура приложения, управляющая жизненным циклом всех компонентов.
type App struct {
	diContainer *diContainer

	httpServer *http.Server
}

// New создает и инициализирует приложение
func New(ctx context.Context) *App {
	a := &App{}

	a.initDeps(ctx)

	return a
}

// Run управляет жизненным циклом приложения:
// - запускает сервер
// - обрабатывает сигналы ОС
// - выполняет graceful shutdown

// Сервер запускается в отдельной горутине, а main-горутина синхронно ждет
// либо сигнал SIGINT/SIGTERM, либо падение сервера. После этого
// closer.CloseAll вызывается синхронно - main-горутина гарантированно
// дожидается завершения всех закрытий перед выходом из Run
func (a *App) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	a.startGracefulShutdown(ctx, cancel)

	return a.runHTTPServer()
}

func (a *App) initDeps(ctx context.Context) {
	inits := []func(context.Context){
		a.initDI,
		a.initLogger,
		a.initHTTPServer,
	}

	for _, f := range inits {
		f(ctx)
	}
}

// initDI создает DI-контейнер.
func (a *App) initDI(_ context.Context) {
	a.diContainer = &diContainer{}
}

// initLogger настраивает глобальный slog с уровнем из конфига.
func (a *App) initLogger(_ context.Context) {
	logger.Init(config.AppConfig().Logger.Level)
}

func (a *App) initHTTPServer(ctx context.Context) {
	server := &http.Server{
		Addr:              config.AppConfig().HTTP.Address(),
		Handler:           a.diContainer.OrderServer(ctx),
		ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атаки
		ReadTimeout:       readTimeout,       // Лимит на чтение всего запроса
		WriteTimeout:      writeTimeout,      // Лимит на запись ответа
		IdleTimeout:       idleTimeout,       // Таймаут keep-alive соединения
		MaxHeaderBytes:    maxHeaderBytes,    // 1 MB // Максимальный размер заголовка
	}

	a.httpServer = server

	closer.Add("http server", func(ctx context.Context) error {
		if err := a.httpServer.Shutdown(ctx); err != nil {
			slog.Error("ошибка при остановке сервера", "error", err)
		}
		return nil
	})
}

func (a *App) startGracefulShutdown(ctx context.Context, cancel context.CancelFunc) {
	go func() {
		<-ctx.Done()

		cancel()

		slog.Info("получен сигнал завершения, начинаем graceful shutdown")

		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, shutdownTimeout)
		defer shutdownCancel()

		if closeErr := closer.CloseAll(shutdownCtx); closeErr != nil {
			slog.Error("ошибка при завершении работы", "error", closeErr)
		}
	}()
}

// runHTTPServer запускает HTTP-сервер и блокирует его до остановки.
func (a *App) runHTTPServer() error {
	// example Graceful shutdown & race condition

	slog.Info("http-сервер запущен", "address", config.AppConfig().HTTP.Address())

	listenErr := a.httpServer.ListenAndServe()
	if listenErr != nil && !errors.Is(listenErr, http.ErrServerClosed) {
		slog.Error("ошибка запуска сервера", "error", listenErr)
	}
	return listenErr
}
