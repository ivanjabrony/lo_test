package server

import (
	"ivanjabrony/test_lo/internal/config"
	"ivanjabrony/test_lo/internal/handler"
	"ivanjabrony/test_lo/internal/middleware"
	"net/http"
)

type Logger interface {
	Log(format string, info ...any)
}

func NewHTTP(
	cfg *config.Config,
	logger Logger,
	taskHandler *handler.TaskHandler) (*http.Server, error) {

	r := http.NewServeMux()

	r.HandleFunc("GET /tasks", taskHandler.HandleGetAllTasks)
	r.HandleFunc("GET /tasks/{task_id}", taskHandler.HandleGetTaskById)
	r.HandleFunc("POST /tasks", taskHandler.HandlePostTask)
	r.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mw := middleware.NewLoggerMiddleware(logger)

	return &http.Server{
		Addr:    ":" + cfg.HttpPort,
		Handler: mw.Logging(r),
	}, nil
}
