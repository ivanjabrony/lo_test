package app

import (
	"context"
	"errors"
	"io"
	"ivanjabrony/test_lo/internal/config"
	"ivanjabrony/test_lo/internal/handler"
	"ivanjabrony/test_lo/internal/storage"
	"ivanjabrony/test_lo/internal/usecase"
	"ivanjabrony/test_lo/pkg/logger"
	"os"
)

type Logger interface {
	Log(format string, info ...any)
}

func InitializeLogger(ctx context.Context, cfg *config.Config, w io.Writer) (Logger, error) {
	if w == nil {
		w = os.Stdout
	}
	logger := logger.NewAsync(ctx, w)
	return logger, nil
}

func InitializeAdapters(cfg *config.Config, logger Logger) (*Handlers, error) {
	if cfg == nil {
		return nil, errors.New("nil values in constructor")
	}

	storages, err := initStorages(logger)
	if err != nil {
		return nil, err
	}

	usecases, err := initUsecases(storages, logger)
	if err != nil {
		return nil, err
	}

	handlers, err := initHandlers(usecases, logger)
	if err != nil {
		return nil, err
	}

	return handlers, nil
}

type Storages struct {
	Task *storage.TaskStorage
}

type Usecases struct {
	Task *usecase.TaskUsecase
}

type Handlers struct {
	Task *handler.TaskHandler
}

func initStorages(logger Logger) (*Storages, error) {
	taslRepository, err := storage.NewTaskStorage(logger)
	if err != nil {
		return nil, err
	}

	return &Storages{
		Task: taslRepository,
	}, nil
}

func initUsecases(storages *Storages, logger Logger) (*Usecases, error) {
	taskUsecase, err := usecase.NewTaskUsecase(logger, storages.Task)
	if err != nil {
		return nil, err
	}

	return &Usecases{
		Task: taskUsecase,
	}, nil
}

func initHandlers(usecases *Usecases, logger Logger) (*Handlers, error) {
	taskHandler, err := handler.NewTaskHandler(logger, usecases.Task)
	if err != nil {
		return nil, err
	}
	return &Handlers{taskHandler}, nil
}
