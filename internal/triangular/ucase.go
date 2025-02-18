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

func New(cfg *appctx.Config, spot providers.Exchange, arbitrage arbitrage.Resolverer) Resolverer {
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
			for _, pair := range t.cfg.TriangularBinancePairs {
				func(p appctx.TriangularBinancePair) {
					g.Go(func() error {
						resp, err := t.arbitrage.Resolve(consts.Binance).PriceByTradingPair(gCtx, p, data)

						if err != nil {
							return fmt.Errorf("get price by trading pair error: %v", err)
						}

						if resp == nil {
							return nil
						}

						rsp, err := t.arbitrage.Resolve(consts.Binance).Calculate(gCtx, *resp)

						if err != nil {
							return fmt.Errorf("calculate error: %v", err)
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
			for _, pair := range t.cfg.TriangularKucoinPairs {
				func(p appctx.TriangularKucoinPair) {
					g.Go(func() error {
						resp, err := t.arbitrage.Resolve(consts.Kucoin).PriceByTradingPair(gCtx, p, data)

						if err != nil {
							return fmt.Errorf("get price by trading pair error: %v", err)
						}

						if resp == nil {
							return nil
						}

						rsp, err := t.arbitrage.Resolve(consts.Kucoin).Calculate(gCtx, *resp)

						if err != nil {
							return fmt.Errorf("calculate error: %v", err)
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
			logger.ErrorWithContext(ctx, fmt.Sprintf("error: %v", err), lf...)
			return err
		}

		if len(tradePairs) > 0 {

			sort.Slice(tradePairs, func(i, j int) bool {
				return tradePairs[i].FinalFunds > tradePairs[j].FinalFunds
			})

			fmt.Println(util.ToJSON(tradePairs[0]))

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

			trade := []providers.PlaceOrderRequest{
				{
					Symbol:   tradePairs[0].PairA,
					Side:     consts.OrderSideBuy,
					Type:     consts.OrderTypeMarket,
					Quantity: tradePairs[0].QtyPairA,
				},
				{
					Symbol: tradePairs[0].PairB,
					Side: func() string {
						if tradePairs[0].Direction == consts.DirectionForward {
							return consts.OrderSideBuy
						}
						return consts.OrderSideSell
					}(),
					Type:     consts.OrderTypeMarket,
					Quantity: tradePairs[0].QtyPairB,
				},
				{
					Symbol:   tradePairs[0].PairC,
					Side:     consts.OrderSideSell,
					Type:     consts.OrderTypeMarket,
					Quantity: tradePairs[0].QtyPairC,
				},
			}

			for _, order := range trade {
				respOrder, err := t.spot.SetExchange(exchange).PlaceOrder(ctx, order)

				lf.Append(logger.Any("raw_request", util.ToJSON(order)))
				lf.Append(logger.Any("raw_response", respOrder.RawResponse()))
				lf.Append(logger.Any("status_code", respOrder.Code))

				if err != nil {
					logger.ErrorWithContext(ctx, fmt.Sprintf("place order error: %v", err), lf...)
					return err
				}

				logger.InfoWithContext(ctx, "success place order", lf...)
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
