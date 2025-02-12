package triangular

import (
	"context"
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

func New(cfg *appctx.Config, spot providers.Exchange, arbitrage arbitrage.Resolverer) Resolve {
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
			logger.EventName("ucase.triangular.start"),
		)
	)

	for _, exchange := range t.cfg.Triangular.Exchanges {

		if !t.cfg.Binance.TriangularEnabled {
			continue
		}

		if !t.cfg.Kucoin.TriangularEnabled {
			continue
		}

		resp, err := t.spot.SetExchange(exchange).TickerPrices(ctx)

		lf.Append(logger.Any("exchange", exchange))
		if err != nil {
			lf.Append(logger.Any("raw_response", resp.RawResponse()))
			lf.Append(logger.Any("status_code", resp.Code))
			logger.ErrorWithContext(ctx, fmt.Sprintf("get ticker price pairs error: %v", err), lf...)
			return err
		}

		data, ok := resp.Data.([]providers.TickerPrices)
		if !ok {
			logger.ErrorWithContext(ctx, "invalid data type", lf...)
			return fmt.Errorf("invalid data type")
		}

		g, gCtx := errgroup.WithContext(ctx)
		if exchange == consts.Binance {
			for _, pair := range t.cfg.TriangularPairs {
				func(p appctx.TriangularBinancePair) {
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

		if exchange == consts.Kucoin {
			for _, pair := range t.cfg.TriangularPairs {
				func(p appctx.TriangularBinancePair) {
					g.Go(func() error {
						resp, err := t.arbitrage.Resolve(consts.Kucoin).PriceByTradingPair(gCtx, p, data)

						if err != nil {
							return err
						}

						rsp, err := t.arbitrage.Resolve(consts.Kucoin).Calculate(gCtx, p, resp)

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

		if err := g.Wait(); err != nil {
			return err
		}

		if len(tradePairs) > 0 {

			if exchange == consts.Binance {
				if !t.cfg.Binance.TradeEnabled {
					continue
				}
				t.cfg.Binance.TriangularEnabled = false
			}

			if exchange == consts.Kucoin {
				if !t.cfg.Kucoin.TradeEnabled {
					continue
				}
				t.cfg.Kucoin.TriangularEnabled = false

			}

			sort.Slice(tradePairs, func(i, j int) bool {
				return tradePairs[i].FinalBalance > tradePairs[j].FinalBalance
			})

			fmt.Println(util.ToJSON(tradePairs))

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
				respOrder, err := t.spot.SetExchange(exchange).PlaceOrder(ctx, order)
				if err != nil {
					return fmt.Errorf("%s place order error: %v, request %v, raw_response: %v, status_code: %v", exchange, err, order, respOrder.RawResponse(), respOrder.Code)
				}
				logger.InfoWithContext(ctx, fmt.Sprintf("%s success place order, raw_request: %v", exchange, util.ToJSON(order)), lf...)
			}
		}

		if exchange == consts.Binance {
			t.cfg.Binance.TriangularEnabled = true
		}

		if exchange == consts.Kucoin {
			t.cfg.Kucoin.TriangularEnabled = true
		}
	}

	return nil
}
