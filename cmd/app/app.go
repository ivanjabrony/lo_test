package app

import (
	"context"
	"errors"
	"fmt"
	"ivanjabrony/test_lo/internal/config"
	"ivanjabrony/test_lo/internal/server"
	"log"
	"net/http"
	"os"
	"time"
)

type Application struct {
	cfg  *config.Config
	http *http.Server
	ctx  context.Context
}

func NewApplication(cfg *config.Config) (*Application, error) {
	if cfg == nil {
		return nil, errors.New("config must be non nil")
	}

	ctx := context.Background()
	logger, err := InitializeLogger(ctx, cfg, os.Stdout)
	if err != nil {
		return nil, err
	}

	handlers, err := InitializeAdapters(cfg, logger)
	if err != nil {
		return nil, err
	}

	http, err := server.NewHTTP(cfg, logger, handlers.Task)
	if err != nil {
		return nil, err
	}

	app := Application{
		cfg:  cfg,
		http: http,
		ctx:  ctx,
	}

	return &app, nil
}

func (app *Application) GetAddr() string {
	return app.http.Addr
}

func (app *Application) Run() error {
	log.Print("Running application")

	ctx, cancel := context.WithCancel(app.ctx)
	defer cancel()

	log.Printf("Starting HTTP server at port: %s", app.cfg.HttpPort)

	serverErr := make(chan error, 1)
	go func() {
		if err := app.http.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
			cancel()
		}
	}()

	select {
	case err := <-serverErr:
		return fmt.Errorf("HTTP server error: %w", err)
	case <-ctx.Done():
		return nil
	}
}

func (app *Application) Stop() {
	log.Print("Shutting down application...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.http.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	log.Print("Application stopped gracefully")
}
