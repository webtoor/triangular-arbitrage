package providers

import (
	"context"

	"github.com/webtoor/triangular-arbitrage/internal/appctx"
)

type ExchangeList map[string]SpotAPI

type Exchange interface {
	SetExchange(name string) SpotAPI
}

type SpotAPI interface {
	GetAllSymbols(ctx context.Context) (appctx.Response, error)
	TickerPrices(ctx context.Context) (appctx.Response, error)
	PlaceOrder(ctx context.Context, in any) (appctx.Response, error)
	GetOrderByID(ctx context.Context, id, symbol string) (appctx.Response, error)
}
