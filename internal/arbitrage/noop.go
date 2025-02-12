package arbitrage

import (
	"context"
	"fmt"

	"github.com/webtoor/triangular-arbitrage/internal/appctx"
	"github.com/webtoor/triangular-arbitrage/internal/providers"
)

type noop struct {
}

func (n *noop) GenerateTriangularPairs(ctx context.Context) error {
	return fmt.Errorf("invalid exchange")
}

func (n *noop) PriceByTradingPair(ctx context.Context, pair appctx.TriangularBinancePair, in []providers.TickerPrices) (PriceByTradingPairResp, error) {
	return PriceByTradingPairResp{}, fmt.Errorf("invalid exchange")
}

func (n *noop) Calculate(ctx context.Context, pair appctx.TriangularBinancePair, prices PriceByTradingPairResp) (*TriangularTradeParam, error) {
	return &TriangularTradeParam{}, fmt.Errorf("invalid exchange")
}
