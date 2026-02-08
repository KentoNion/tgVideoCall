package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"tgVideoCall/app"
	"tgVideoCall/pkg/config"
	"tgVideoCall/pkg/logger"
)

func main() {
	cfg := config.MustLoad()
	log := logger.MustInitLogger(cfg)
	ctx, cancel := context.WithCancel(context.Background())

	app.Run(ctx, log, cfg)
	go shutdown(cancel, log)
}

func shutdown(cancel context.CancelFunc, log slog.Logger) {
	const op = "shutdown"	

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	// Ожидание сигнала
	receivedSignal := <-sigterm
	log.Info(op, "recieved signal", receivedSignal)
	cancel()
}
