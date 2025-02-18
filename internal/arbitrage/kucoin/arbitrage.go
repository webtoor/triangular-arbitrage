package arbitragekucoin

import (
	"context"
	"fmt"

	"github.com/spf13/cast"
	"github.com/webtoor/triangular-arbitrage/internal/appctx"
	"github.com/webtoor/triangular-arbitrage/internal/arbitrage"
	"github.com/webtoor/triangular-arbitrage/internal/common"
	"github.com/webtoor/triangular-arbitrage/internal/consts"
	"github.com/webtoor/triangular-arbitrage/internal/presentations"
	"github.com/webtoor/triangular-arbitrage/internal/providers"
	"github.com/webtoor/triangular-arbitrage/pkg/file"
	"github.com/webtoor/triangular-arbitrage/pkg/util"
)

type kucoin struct {
	cfg  *appctx.Config
	spot providers.Exchange
}

func New(cfg *appctx.Config, spot providers.Exchange) arbitrage.Triangular {
	return &kucoin{
		cfg:  cfg,
		spot: spot,
	}
}

func (a *kucoin) Calculate(ctx context.Context, prices arbitrage.PriceByTradingPairResp) (*arbitrage.TriangularTradeParam, error) {

	var (
		resp       arbitrage.TriangularTradeParam
		directions = []string{consts.DirectionForward, consts.DirectionReverse}
	)

	if prices.PairAAsk == 0 || prices.PairABid == 0 || prices.PairBAsk == 0 || prices.PairBBid == 0 || prices.PairCAsk == 0 || prices.PairCBid == 0 {
		return nil, nil
	}

	for _, direction := range directions {

		if direction == consts.DirectionForward {
			// USDT -> BTC
			tradePairA := a.cfg.Triangular.Balance / prices.PairAAsk
			// USDT -> BTC
			tradeWithFeePairA := (a.cfg.Triangular.Balance / prices.PairAAsk) - ((a.cfg.Triangular.Balance / prices.PairAAsk) * a.cfg.Binance.Fees)
			// BTC -> ETH
			tradeWithFeePairB := (tradeWithFeePairA / prices.PairBAsk) - ((tradeWithFeePairA / prices.PairBAsk) * a.cfg.Binance.Fees)
			tradeWithFeePairB2 := tradeWithFeePairB - (tradeWithFeePairB * a.cfg.Binance.Fees)
			// ETH -> USDT
			tradeWithFeePairC := (tradeWithFeePairB * prices.PairCBid) - ((tradeWithFeePairB * prices.PairCBid) * a.cfg.Binance.Fees)

			// Filter Profitable
			if tradeWithFeePairC > (a.cfg.Triangular.Balance + (a.cfg.Triangular.Balance * a.cfg.Triangular.Profit)) {
				// Filter Quantity
				if cast.ToFloat64(prices.PairAAskQty) > tradePairA && cast.ToFloat64(prices.PairBAskQty) > tradeWithFeePairB && cast.ToFloat64(prices.PairCBidQty) > tradeWithFeePairB2 {
					resp.Exchange = consts.Kucoin
					resp.Direction = consts.DirectionForward
					resp.PairA = prices.PairA
					resp.PairB = prices.PairB
					resp.PairC = prices.PairC
					resp.PricesA = prices.PairAAsk
					resp.PricesB = prices.PairBAsk
					resp.PricesC = prices.PairABid
					resp.QtyPairA = util.Precision(tradePairA, prices.PairAAskQty)
					resp.QtyPairB = util.Precision(tradeWithFeePairB, prices.PairBAskQty)
					resp.QtyPairC = util.Precision(tradeWithFeePairB2, prices.PairCBidQty)
					resp.InitialFunds = a.cfg.Triangular.Balance
					resp.FinalFunds = util.Round(tradeWithFeePairC, 4)
					return &resp, nil
				}
			}
		}

		if direction == consts.DirectionReverse {
			// USDT -> ETH
			tradePairA := a.cfg.Triangular.Balance / prices.PairCAsk
			// USDT -> ETH
			tradeWithFeePairA := (a.cfg.Triangular.Balance / prices.PairCAsk) - ((a.cfg.Triangular.Balance / prices.PairCAsk) * a.cfg.Binance.Fees)
			// ETH -> BTC
			tradeWithFeePairB := (tradeWithFeePairA * prices.PairBBid) - ((tradeWithFeePairA * prices.PairBBid) * a.cfg.Binance.Fees)
			// BTC -> USDT
			tradeWithFeePairC := (tradeWithFeePairB * prices.PairABid) - ((tradeWithFeePairB * prices.PairABid) * a.cfg.Binance.Fees)

			// Filter Profitable
			if tradeWithFeePairC > (a.cfg.Triangular.Balance + (a.cfg.Triangular.Balance * a.cfg.Triangular.Profit)) {
				// Filter Quantity
				if cast.ToFloat64(prices.PairCAskQty) > tradePairA && cast.ToFloat64(prices.PairBBidQty) > tradeWithFeePairA && cast.ToFloat64(prices.PairABidQty) > tradeWithFeePairB {
					resp.Exchange = consts.Kucoin
					resp.Direction = consts.DirectionReverse
					resp.PairA = prices.PairC
					resp.PairB = prices.PairB
					resp.PairC = prices.PairA
					resp.PricesA = prices.PairCAsk
					resp.PricesB = prices.PairBBid
					resp.PricesC = prices.PairABid
					resp.QtyPairA = util.Precision(tradePairA, prices.PairCAskQty)
					resp.QtyPairB = util.Precision(tradeWithFeePairA, prices.PairBBidQty)
					resp.QtyPairC = util.Precision(tradeWithFeePairB, prices.PairABidQty)
					resp.InitialFunds = a.cfg.Triangular.Balance
					resp.FinalFunds = util.Round(tradeWithFeePairC, 4)
					return &resp, nil
				}
			}
		}
	}

	return nil, nil
}

func (a *kucoin) PriceByTradingPair(ctx context.Context, pair any, in []providers.TickerPrices) (*arbitrage.PriceByTradingPairResp, error) {

	var rsp arbitrage.PriceByTradingPairResp

	data, ok := pair.(appctx.TriangularKucoinPair)
	if !ok {
		return &rsp, fmt.Errorf("invalid type triangular pair")
	}

	rsp.PairA = data.PairA
	rsp.PairB = data.PairB
	rsp.PairC = data.PairC
	rsp.PairAAsk, rsp.PairABid, rsp.PairAAskQty, rsp.PairABidQty = arbitrage.ExtractPrice(data.PairA, in)
	rsp.PairBAsk, rsp.PairBBid, rsp.PairBAskQty, rsp.PairBBidQty = arbitrage.ExtractPrice(data.PairB, in)
	rsp.PairCAsk, rsp.PairCBid, rsp.PairCAskQty, rsp.PairCBidQty = arbitrage.ExtractPrice(data.PairC, in)

	return &rsp, nil
}

func (a *kucoin) GenerateTriangularPairs(ctx context.Context) error {
	var (
		keys  = make(map[string]bool)
		pairs = []presentations.TriangularCombined{}
	)

	resp, err := a.spot.SetExchange(consts.Kucoin).GetAllSymbols(ctx)

	if err != nil {
		return fmt.Errorf("generate triangular pairs kucoin error: %v, raw response: %v, status code: %v", err, resp.RawResponse(), resp.Code)
	}

	symbols := common.TradableSymbols(resp.Data)

	pair_a := fmt.Sprintf("%s-%s", consts.ABaseBTC, consts.AQuoteUSDT)
	a_base := consts.ABaseBTC
	a_quote := consts.AQuoteUSDT

	for range symbols {
		a_pair_box := []string{a_base, a_quote}

		for _, pair_b := range symbols {
			if pair_b["quote"] == a_base {
				b_base := pair_b["base"]
				b_quote := pair_b["quote"]

				if pair_b["symbol"] != pair_a {
					if util.InArray(b_base, a_pair_box) || util.InArray(b_quote, a_pair_box) {
						for _, pair_c := range symbols {
							if pair_c["quote"] == a_quote {
								c_base := pair_c["base"]
								c_quote := pair_c["quote"]
								if pair_c["symbol"] != pair_a && pair_c["symbol"] != pair_b["symbol"] {
									pair_box := []string{a_base, a_quote, b_base, b_quote, c_base, c_quote}

									counts_c_base := 0
									for _, v := range pair_box {
										if v == c_base {
											counts_c_base++
										}
									}

									counts_c_quote := 0
									for _, v := range pair_box {
										if v == c_base {
											counts_c_quote++
										}
									}

									if counts_c_base == 2 && counts_c_quote == 2 && c_base != c_quote {
										combined := fmt.Sprintf("%s,%s,%s", pair_a, pair_b["symbol"], pair_c["symbol"])
										if _, v := keys[combined]; !v {
											keys[combined] = true
											pairs = append(pairs, presentations.TriangularCombined{
												ABase:    a_base,
												AQuote:   a_quote,
												BBase:    b_base,
												BQuote:   b_quote,
												CBase:    c_base,
												CQuote:   c_quote,
												PairA:    pair_a,
												PairB:    pair_b["symbol"],
												PairC:    pair_c["symbol"],
												Combined: combined,
											})
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	err = file.WriteToJson(a.cfg.Kucoin.PathTriangularPairs, pairs)

	if err != nil {
		return err
	}

	return nil
}
