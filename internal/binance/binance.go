package binance

import (
	"context"
	"finbin/internal/logger"
	"github.com/aiviaio/go-binance/v2"
)

const TradingStatus = "TRADING"

type ClientWrapper struct {
	client *binance.Client
	logger *logger.Logger
}

func NewBinanceClientWrapper(apiKey, secretKey string, log *logger.Logger) *ClientWrapper {
	client := binance.NewClient(apiKey, secretKey)
	return &ClientWrapper{
		client: client,
		logger: log,
	}
}

// FetchFirstNSymbols retrieves a limited number of trading symbols from the exchange.
func (b *ClientWrapper) FetchFirstNSymbols(ctx context.Context, symbolsLimit int) ([]string, error) {
	exchangeInfo, err := b.client.NewExchangeInfoService().Do(ctx)
	if err != nil {
		b.logger.Errorf("Error fetching exchange info: %v", err)
		return nil, err
	}

	return filterTradingSymbols(exchangeInfo.Symbols, symbolsLimit), nil
}

// filterTradingSymbols filters out symbols that are currently trading.
func filterTradingSymbols(symbols []binance.Symbol, limit int) []string {
	tradingSymbols := make([]string, 0, limit)
	for _, symbol := range symbols {
		if symbol.Status == TradingStatus {
			tradingSymbols = append(tradingSymbols, symbol.Symbol)
			if len(tradingSymbols) == limit {
				break
			}
		}
	}
	return tradingSymbols
}

// FetchPrice retrieves the current price of a symbol.
func (b *ClientWrapper) FetchPrice(ctx context.Context, symbol string) (string, error) {
	prices, err := b.client.NewListPricesService().Symbol(symbol).Do(ctx)
	if err != nil {
		b.logger.Errorf("Error fetching price for symbol %s: %v", symbol, err)
		return "", err
	}

	if len(prices) == 0 {
		b.logger.Errorf("No price data for symbol %s", symbol)
		return "", nil
	}

	return prices[0].Price, nil
}
