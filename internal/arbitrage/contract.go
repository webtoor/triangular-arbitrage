package arbitrage

import (
	"context"

	"github.com/webtoor/triangular-arbitrage/internal/providers"
)

type ArbitrageList map[string]Triangular

type Resolverer interface {
	Resolve(name string) Triangular
}

type Triangular interface {
	GenerateTriangularPairs(ctx context.Context) error
	PriceByTradingPair(ctx context.Context, pair any, in []providers.TickerPrices) (*PriceByTradingPairResp, error)
	Calculate(ctx context.Context, prices PriceByTradingPairResp) (*TriangularTradeParam, error)
}
