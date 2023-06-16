// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"securewebapp/config"
	v1 "securewebapp/internal/controller/http/v1"
	"securewebapp/internal/usecase"
	"securewebapp/internal/usecase/repo"
	"securewebapp/pkg/httpserver"
	"securewebapp/pkg/logger"
	"securewebapp/pkg/postgres"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Use case
	userUseCase := usecase.New(repo.New(pg))

	// HTTP Server
	handler := gin.New()

	v1.NewRouter(handler, l, *userUseCase, cfg.CsrfSecret, cfg.CookieSecret)
	httpServer := httpserver.New(handler, cfg.HTTP.TlsCert, cfg.HTTP.TlsKey, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

}
