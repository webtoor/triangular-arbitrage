package arbitragebinance

import (
	"context"
	"fmt"

	"github.com/webtoor/triangular-arbitrage/internal/appctx"
	"github.com/webtoor/triangular-arbitrage/internal/arbitrage"
	"github.com/webtoor/triangular-arbitrage/internal/common"
	"github.com/webtoor/triangular-arbitrage/internal/consts"
	"github.com/webtoor/triangular-arbitrage/internal/presentations"
	"github.com/webtoor/triangular-arbitrage/internal/providers"
	"github.com/webtoor/triangular-arbitrage/pkg/file"
	"github.com/webtoor/triangular-arbitrage/pkg/util"
)

type binance struct {
	cfg  *appctx.Config
	spot providers.Exchange
}

func New(cfg *appctx.Config, spot providers.Exchange) arbitrage.Triangular {
	return &binance{
		cfg:  cfg,
		spot: spot,
	}
}

func (t *binance) Calculate(ctx context.Context, prices arbitrage.PriceByTradingPairResp) (*arbitrage.TriangularTradeParam, error) {

	var (
		resp       arbitrage.TriangularTradeParam
		directions = []string{consts.DirectionForward, consts.DirectionReverse}
	)

	if prices.PairAAsk == 0 || prices.PairABid == 0 || prices.PairBAsk == 0 || prices.PairBBid == 0 || prices.PairCAsk == 0 || prices.PairCBid == 0 {
		return nil, nil
	}

	for _, direction := range directions {
		if direction == consts.DirectionForward {
			tradePairA := t.cfg.Triangular.Balance / prices.PairAAsk
			tradePairB := tradePairA / prices.PairBAsk
			_ = tradePairB * prices.PairCBid

			tradeWithFeePairA := (t.cfg.Triangular.Balance / prices.PairAAsk) - ((t.cfg.Triangular.Balance / prices.PairAAsk) * t.cfg.Binance.Fees)
			tradeWithFeePairB := (tradeWithFeePairA / prices.PairBAsk) - ((tradeWithFeePairA / prices.PairBAsk) * t.cfg.Binance.Fees)
			tradeWithFeePairC := (tradeWithFeePairB * prices.PairCBid) - ((tradeWithFeePairB * prices.PairCBid) * t.cfg.Binance.Fees)

			if tradeWithFeePairC > t.cfg.Triangular.Balance {
				resp.Exchange = consts.Binance
				resp.Direction = consts.DirectionForward
				resp.PairA = prices.PairA
				resp.PairB = prices.PairB
				resp.PairC = prices.PairC
				resp.QtyPairA = t.cfg.Triangular.Balance
				resp.QtyPairB = tradeWithFeePairA
				resp.QtyPairC = tradeWithFeePairC
				resp.InitialFunds = t.cfg.Triangular.Balance
				resp.FinalFunds = tradeWithFeePairC
				return &resp, nil
			}
		}

		if direction == consts.DirectionReverse {
			tradePairA := t.cfg.Triangular.Balance / prices.PairCAsk
			tradePairB := tradePairA * prices.PairBBid
			_ = tradePairB * prices.PairABid

			tradeWithFeePairA := (t.cfg.Triangular.Balance / prices.PairCAsk) - ((t.cfg.Triangular.Balance / prices.PairCAsk) * t.cfg.Binance.Fees)
			tradeWithFeePairB := (tradeWithFeePairA * prices.PairBBid) - ((tradeWithFeePairA * prices.PairBBid) * t.cfg.Binance.Fees)
			tradeWithFeePairC := (tradeWithFeePairB * prices.PairABid) - ((tradeWithFeePairB * prices.PairABid) * t.cfg.Binance.Fees)

			if tradeWithFeePairC > t.cfg.Triangular.Balance {
				resp.Exchange = consts.Binance
				resp.Direction = consts.DirectionReverse
				resp.PairA = prices.PairC
				resp.PairB = prices.PairB
				resp.PairC = prices.PairA
				resp.QtyPairA = t.cfg.Triangular.Balance
				resp.QtyPairB = tradeWithFeePairB
				resp.QtyPairC = tradeWithFeePairC
				resp.InitialFunds = t.cfg.Triangular.Balance
				resp.FinalFunds = tradeWithFeePairC
				return &resp, nil
			}
		}
	}

	return nil, nil
}

func (t *binance) PriceByTradingPair(ctx context.Context, pair any, in []providers.TickerPrices) (arbitrage.PriceByTradingPairResp, error) {

	rsp := arbitrage.PriceByTradingPairResp{}

	data, ok := pair.(appctx.TriangularBinancePair)
	if !ok {
		return rsp, fmt.Errorf("invalid type triangular pair")
	}

	rsp.PairA = data.PairA
	rsp.PairB = data.PairB
	rsp.PairC = data.PairC
	rsp.PairAAsk, rsp.PairABid = arbitrage.ExtractPrice(data.PairA, in)
	rsp.PairBAsk, rsp.PairBBid = arbitrage.ExtractPrice(data.PairB, in)
	rsp.PairCAsk, rsp.PairCBid = arbitrage.ExtractPrice(data.PairC, in)

	return rsp, nil
}

func (t *binance) GenerateTriangularPairs(ctx context.Context) error {
	var (
		keys  = make(map[string]bool)
		pairs = []presentations.TriangularCombined{}
	)

	resp, err := t.spot.SetExchange(consts.Binance).GetAllSymbols(ctx)

	if err != nil {
		return fmt.Errorf("generate triangular pairs binance error: %v, raw response: %v, status code: %v", err, resp.RawResponse(), resp.Code)
	}

	symbols := common.TradableSymbols(resp.Data)

	pair_a := fmt.Sprintf("%s%s", consts.ABaseBTC, consts.AQuoteUSDT)
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

	err = file.WriteToJson(t.cfg.Binance.PathTriangularPairs, pairs)

	if err != nil {
		return err
	}

	return nil
}
