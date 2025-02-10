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
	ExchangeInfo(ctx context.Context) (appctx.Response, error)
	TickerPrices(ctx context.Context) (appctx.Response, error)
}
