package providers

import (
	"context"
	"fmt"

	"github.com/webtoor/triangular-arbitrage/internal/appctx"
)

type noop struct {
}

func (n *noop) ExchangeInfo(ctx context.Context) (appctx.Response, error) {
	return *appctx.NewResponse(), fmt.Errorf("invalid exchange")
}

func (n *noop) TickerPrices(ctx context.Context) (appctx.Response, error) {
	return *appctx.NewResponse(), fmt.Errorf("invalid exchange")
}
