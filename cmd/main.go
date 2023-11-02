package main

import (
	"context"
	"finbin/internal/app"
	"finbin/internal/binance"
	"finbin/internal/config"
	"finbin/internal/logger"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log := initializeLogger()
	cfg := loadConfiguration(log)
	binanceClient := initializeBinanceClient(cfg, log)
	myApp := app.New(cfg, log, binanceClient)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	runApplication(ctx, myApp, log)

	initiateShutdown(cfg, myApp, log)
}

func initializeLogger() *logger.Logger {
	log, err := logger.New()
	if err != nil {
		fmt.Printf("Error initializing logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()
	return log
}

func loadConfiguration(log *logger.Logger) *config.Config {
	cfg, err := config.NewConfig()
	handleErrors(err, log, "Error loading configuration: %v")
	log.Infof("Loaded Configuration: %+v", cfg)
	return cfg
}

func initializeBinanceClient(cfg *config.Config, log *logger.Logger) *binance.ClientWrapper {
	binanceClient := binance.NewBinanceClientWrapper(cfg.BinanceAPIKey, cfg.BinanceSecret, log)
	log.Info("Binance Client initialized")
	return binanceClient
}

func runApplication(ctx context.Context, myApp *app.App, log *logger.Logger) {
	errChan := make(chan error, 1)
	go func() {
		errChan <- myApp.Run(ctx)
	}()

	select {
	case <-ctx.Done():
		log.Info("Shutdown signal received from OS.")
	case err := <-errChan:
		handleErrors(err, log, "Application run terminated: %v")
	}
}

func initiateShutdown(cfg *config.Config, myApp *app.App, log *logger.Logger) {
	log.Info("Starting shutdown process...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.GracefulTimeout)
	defer cancel()

	err := myApp.Shutdown(shutdownCtx)
	handleErrors(err, log, "Error during shutdown: %v")
}

func handleErrors(err error, log *logger.Logger, message string) {
	if err != nil {
		log.Errorf(message, err)
		os.Exit(1)
	}
}
