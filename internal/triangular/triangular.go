package triangular

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/webtoor/triangular-arbitrage/internal/appctx"
	"github.com/webtoor/triangular-arbitrage/internal/arbitrage"
	"github.com/webtoor/triangular-arbitrage/internal/consts"
	"github.com/webtoor/triangular-arbitrage/internal/providers"
	"github.com/webtoor/triangular-arbitrage/pkg/logger"
	"github.com/webtoor/triangular-arbitrage/pkg/util"
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
		lf         = logger.NewFields(
			logger.EventName("Triangular.Start"),
		)
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

	if len(tradePairs) > 0 {
		t.cfg.Binance.TriangularEnabled = false

		sort.Slice(tradePairs, func(i, j int) bool {
			return tradePairs[i].FinalBalance > tradePairs[j].FinalBalance
		})

		x, _ := json.Marshal(tradePairs[0])

		fmt.Println(string(x))

		if t.cfg.Binance.TradeEnabled {
			trade := []providers.PlaceOrderRequest{
				{
					Symbol:   tradePairs[0].PairA,
					Side:     consts.BinanceSideBuy,
					Type:     consts.BinanceTypeMarket,
					Quantity: tradePairs[0].QtyPairA,
				},
				{
					Symbol: tradePairs[0].PairB,
					Side: func() string {
						if tradePairs[0].Direction == consts.DirectionReverse {
							return consts.BinanceSideSell
						}
						return consts.BinanceSideBuy
					}(),
					Type:     consts.BinanceTypeMarket,
					Quantity: util.Round(tradePairs[0].QtyPairB, 8),
				},
				{
					Symbol:   tradePairs[0].PairC,
					Side:     consts.BinanceSideSell,
					Type:     consts.BinanceTypeMarket,
					Quantity: tradePairs[0].QtyPairC,
				},
			}

			for _, order := range trade {
				respOrder, err := t.spot.SetExchange(consts.Binance).PlaceOrder(ctx, order)
				if err != nil {
					return fmt.Errorf("place order error: %v, request %v, raw_response: %v, status_code: %v", err, order, respOrder.RawResponse(), respOrder.Code)
				}
				logger.InfoWithContext(ctx, fmt.Sprintf("success place order, raw_request: %v", util.ToJSON(order)), lf...)
			}
		}
	}

	t.cfg.Binance.TriangularEnabled = true
	return nil
}
