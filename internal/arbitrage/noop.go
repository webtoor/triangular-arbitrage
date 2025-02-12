package arbitrage

import (
	"context"
	"errors"

	"github.com/webtoor/triangular-arbitrage/internal/appctx"
	"github.com/webtoor/triangular-arbitrage/internal/providers"
)

type noop struct {
}

func (n *noop) GenerateTriangularPairs(ctx context.Context) error {
	return errors.New("GenerateTriangularPairs: invalid exchange")
}

func (n *noop) PriceByTradingPair(ctx context.Context, pair appctx.TriangularBinancePair, in []providers.TickerPrices) (PriceByTradingPairResp, error) {
	return PriceByTradingPairResp{}, errors.New("PriceByTradingPair: invalid exchange")
}

func (n *noop) Calculate(ctx context.Context, prices PriceByTradingPairResp) (*TriangularTradeParam, error) {
	return &TriangularTradeParam{}, errors.New("Calculate: invalid exchange")
}
