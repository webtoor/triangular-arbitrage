package triangular

import (
	"context"
	"fmt"
	"sort"

	"github.com/webtoor/triangular-arbitrage/internal/appctx"
	"github.com/webtoor/triangular-arbitrage/internal/arbitrage"
	"github.com/webtoor/triangular-arbitrage/internal/consts"
	"github.com/webtoor/triangular-arbitrage/internal/providers"
	"golang.org/x/sync/errgroup"
)

type triangular struct {
	cfg       *appctx.Config
	spot      providers.Exchange
	arbitrage arbitrage.Resolverer
}

func New(cfg *appctx.Config, spot providers.Exchange, arbitrage arbitrage.Resolverer) *triangular {
	return &triangular{
		cfg:       cfg,
		spot:      spot,
		arbitrage: arbitrage,
	}
}

func (t *triangular) Start(ctx context.Context) error {
	var (
		tradePairs []arbitrage.TriangularTradeParam
	)

	resp, err := t.spot.SetExchange(consts.Binance).TickerPrices(ctx)

	if err != nil {
		return fmt.Errorf("get price pairs binance error: %v, raw response: %v, status code: %v", err, resp.RawResponse(), resp.Code)
	}

	data, ok := resp.Data.([]providers.TickerPrices)
	if !ok {
		return fmt.Errorf("invalid data type")
	}

	g, gCtx := errgroup.WithContext(ctx)
	for _, exchange := range t.cfg.Triangular.Exchanges {
		switch exchange {
		case consts.Binance:
			for _, pair := range t.cfg.TriangularPairs {
				func(p appctx.TriangularPair) {
					g.Go(func() error {
						resp, err := t.arbitrage.Resolve(consts.Binance).PriceByTradingPair(gCtx, p, data)

						if err != nil {
							return err
						}

						rsp, err := t.arbitrage.Resolve(consts.Binance).Calculate(gCtx, p, resp)

						if err != nil {
							return err
						}

						if rsp != nil {
							tradePairs = append(tradePairs, *rsp)
						}

						return err
					})
				}(pair)
			}
		}

	}

	if err := g.Wait(); err != nil {
		return err
	}

	sort.Slice(tradePairs, func(i, j int) bool {
		return tradePairs[i].FinalBalance > tradePairs[j].FinalBalance
	})

	fmt.Println(tradePairs)

	return nil
}
