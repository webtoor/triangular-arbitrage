package common

import (
	"github.com/webtoor/triangular-arbitrage/internal/consts"
	"github.com/webtoor/triangular-arbitrage/internal/providers/binance"
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
	}

	return symbols
}
