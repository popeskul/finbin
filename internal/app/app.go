package app

import (
	"context"
	"finbin/internal/binance"
	"finbin/internal/config"
	"finbin/internal/logger"
	"fmt"
	"sync"
)

type App struct {
	config *config.Config
	logger *logger.Logger
	client *binance.ClientWrapper
	wg     sync.WaitGroup
	stopCh chan struct{}
}

func New(cfg *config.Config, log *logger.Logger, client *binance.ClientWrapper) *App {
	return &App{
		config: cfg,
		logger: log,
		client: client,
		stopCh: make(chan struct{}),
	}
}

func (a *App) Run(ctx context.Context) error {
	a.logger.Info("Starting application components")

	doneCh := make(chan struct{})

	go func() {
		if err := a.startListening(ctx); err != nil {
			a.logger.Errorf("Error occurred during listening: %v", err)
		}
		close(doneCh)
	}()

	select {
	case <-ctx.Done():
		a.logger.Info("Context cancelled, stopping application")
		return ctx.Err()
	case <-doneCh:
		a.logger.Info("startListening finished, stopping application")
		return nil
	}
}

func (a *App) startListening(ctx context.Context) error {
	a.logger.Info("Start listening to data source...")

	symbols, err := a.client.FetchFirstNSymbols(ctx, 5)
	if err != nil {
		a.logger.Errorf("Failed to fetch symbols: %v", err)
		return err
	}
	a.logger.Infof("Fetched symbols: %v", symbols)

	priceCh := make(chan string, 5)
	doneCh := make(chan struct{})

	for _, symbol := range symbols {
		a.wg.Add(1)
		go func(s string) {
			defer a.wg.Done()
			select {
			case <-ctx.Done():
				// If the context is cancelled, exit the goroutine.
				return
			default:
				price, err := a.client.FetchPrice(ctx, s)
				if err != nil {
					a.logger.Errorf("Failed to fetch price for symbol %s: %v", s, err)
					priceCh <- fmt.Sprintf("Error fetching price for symbol %s", s)
					return
				}
				priceCh <- fmt.Sprintf("Price for %s: %s", s, price)
			}
		}(symbol)
	}

	// Start a goroutine to close the priceCh channel once all fetching goroutines are done.
	go func() {
		a.wg.Wait()
		close(priceCh)
		close(doneCh)
	}()

	// Wait for all goroutines to finish and print the prices.
	go func() {
		for p := range priceCh {
			fmt.Println(p)
		}
	}()

	<-doneCh // Wait for the done signal before returning.

	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	a.logger.Info("Application stopped gracefully")
	return nil
}
