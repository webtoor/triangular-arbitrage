package providers

import (
	"context"
	"fmt"

	"github.com/webtoor/triangular-arbitrage/internal/appctx"
)

type noop struct {
}

func (n *noop) GetAllSymbols(ctx context.Context) (appctx.Response, error) {
	return *appctx.NewResponse(), fmt.Errorf("GetAllSymbols: invalid exchange")
}

func (n *noop) TickerPrices(ctx context.Context) (appctx.Response, error) {
	return *appctx.NewResponse(), fmt.Errorf("TickerPrices: invalid exchange")
}

func (n *noop) PlaceOrder(ctx context.Context, in any) (appctx.Response, error) {
	return *appctx.NewResponse(), fmt.Errorf("PlaceOrder: invalid exchange")
}
