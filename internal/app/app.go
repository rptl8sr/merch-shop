package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"merch-shop/internal/config"
	"merch-shop/internal/repository"
	"merch-shop/internal/router"
	"merch-shop/internal/service"
	"merch-shop/pkg/database"
	"merch-shop/pkg/logger"
)

type app struct {
	httpServer *http.Server
	db         database.DB
}

type App interface {
	Run(ctx context.Context) error
	Shutdown() error
}

func Must(ctx context.Context) App {
	a := &app{}

	cfg := config.Must()
	logger.Init(cfg.AppLogLevel)

	db, err := database.New(ctx, cfg.GetDSN(), nil)
	if err != nil {
		logger.Error("database.Must: ", "message", err.Error())
		panic(err)
	}

	repo := repository.NewRepository(db)
	services := service.NewServices(*repo)
	r := router.New(cfg.JWTSecret, services)

	a.db = db
	a.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.AppPort),
		Handler: r,
	}

	return a
}

func (a *app) Run(ctx context.Context) error {
	logger.Info("Starting server", "address", a.httpServer.Addr)

	errChan := make(chan error, 1)

	go func() {
		err := a.httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
		close(errChan)
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return a.Shutdown()
	}
}

func (a *app) Shutdown() error {
	logger.Info("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		logger.Error("Graceful shutdown failed: ", "error", err.Error())
		return err
	}

	a.db.Close()

	logger.Info("Server and database shutdown completed successfully")
	return nil
}
