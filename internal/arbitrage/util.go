package arbitrage

import (
	"github.com/spf13/cast"
	"github.com/webtoor/triangular-arbitrage/internal/providers"
)

func ExtractPrice(pair string, prices []providers.TickerPrices) (ask float64, bid float64) {
	for _, price := range prices {
		if price.Symbol == pair {
			return cast.ToFloat64(price.AskPrice), cast.ToFloat64(price.BidPrice)
		}
	}

	return
}
