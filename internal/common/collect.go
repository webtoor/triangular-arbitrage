package common

import (
	"github.com/webtoor/triangular-arbitrage/internal/consts"
	"github.com/webtoor/triangular-arbitrage/internal/providers/binance"
	"github.com/webtoor/triangular-arbitrage/internal/providers/kucoin"
)

func TradableSymbols(in any) []map[string]string {
	var symbols []map[string]string

	switch v := in.(type) {
	case binance.AllSymbolsResponse:
		for _, symbol := range v.Symbols {
			if symbol.Status == consts.BinanceSymbolStatusTrading {
				symbols = append(symbols, map[string]string{
					"symbol": symbol.Symbol,
					"base":   symbol.BaseAsset,
					"quote":  symbol.QuoteAsset,
				})
			}
		}
	case kucoin.AllSymbolsResponse:
		for _, symbol := range v.Data {
			if symbol.EnableTrading {
				symbols = append(symbols, map[string]string{
					"symbol": symbol.Symbol,
					"base":   symbol.BaseCurrency,
					"quote":  symbol.QuoteCurrency,
				})
			}
		}
	}

	return symbols
}
