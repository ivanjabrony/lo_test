package main

import (
	"ivanjabrony/test_lo/cmd/app"
	"ivanjabrony/test_lo/internal/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	app, err := app.NewApplication(&cfg)
	if err != nil {
		log.Fatalf("Error while starting app: %v", err)
	}

	done := make(chan os.Signal, 2)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	appErr := make(chan error, 2)
	go func() {
		appErr <- app.Run()
	}()

	select {
	case err := <-appErr:
		log.Fatalf("Application error: %v", err)
	case <-done:
		app.Stop()
	}
}
