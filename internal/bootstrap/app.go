package bootstrap

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type App struct {
	Container  *Container
	HTTPServer *fiber.App
}

func NewApp(container *Container, router *fiber.App) *App {
	return &App{
		Container:  container,
		HTTPServer: router,
	}
}

// Start initiates the HTTP server and blocks until a termination signal is received.
func (a *App) Start() error {
	addr := fmt.Sprintf(":%d", a.Container.Config.App.Port)

	// Channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)

	go func() {
		a.Container.Logger.Info("Starting HTTP server", zap.String("addr", addr))
		serverErrors <- a.HTTPServer.Listen(addr)
	}()

	// Channel to listen for interrupt or terminate signals.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		a.Container.Logger.Info("Start shutdown", zap.String("signal", sig.String()))

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := a.HTTPServer.ShutdownWithContext(ctx); err != nil {
			a.Container.Logger.Error("Graceful shutdown did not complete in time", zap.Error(err))
			if err := a.HTTPServer.Shutdown(); err != nil {
				return fmt.Errorf("could not stop server gracefully: %w", err)
			}
		}

		a.Container.Close()
		a.Container.Logger.Info("Shutdown complete")
	}

	return nil
}