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

func (t *binance) Calculate(ctx context.Context, pair appctx.TriangularPair, prices arbitrage.PriceByTradingPairResp) error {

	if prices.PairAAsk == 0 || prices.PairABid == 0 || prices.PairBAAsk == 0 || prices.PairBBid == 0 || prices.PairCAAsk == 0 || prices.PairCBid == 0 {
		return nil
	}

	tradePairA := t.cfg.Triangular.Balance*(1/prices.PairAAsk) - ((t.cfg.Triangular.Balance * (1 / prices.PairAAsk)) * t.cfg.Binance.Fees)
	tradePairB := tradePairA*(1/prices.PairBAAsk) - ((tradePairA * (1 / prices.PairBAAsk)) * t.cfg.Binance.Fees)
	tradePairC := tradePairB*prices.PairCBid - ((tradePairB * prices.PairCBid) * t.cfg.Binance.Fees)

	if tradePairC > t.cfg.Triangular.Balance {
		fmt.Printf("%s:$%v - %s:%v - %s:$%v - startBalance:$%v - endBalance:$%v\n", pair.PairA, prices.PairAAsk, pair.PairB, prices.PairBAAsk, pair.PairC, prices.PairCBid, t.cfg.Triangular.Balance, tradePairC)
	}
	return nil

}

func (t *binance) PriceByTradingPair(ctx context.Context, pair appctx.TriangularPair, in []providers.TickerPrices) (arbitrage.PriceByTradingPairResp, error) {

	rsp := arbitrage.PriceByTradingPairResp{}
	rsp.PairAAsk, rsp.PairABid = arbitrage.ExtractPrice(pair.PairA, in)
	rsp.PairBAAsk, rsp.PairBBid = arbitrage.ExtractPrice(pair.PairB, in)
	rsp.PairCAAsk, rsp.PairCBid = arbitrage.ExtractPrice(pair.PairC, in)

	return rsp, nil
}

func (t *binance) GenerateTriangularPairs(ctx context.Context) error {
	var (
		keys  = make(map[string]bool)
		pairs = []presentations.TriangularCombined{}
	)

	resp, err := t.spot.SetExchange(consts.Binance).ExchangeInfo(ctx)

	if err != nil {
		return fmt.Errorf("generate triangular pairs binance error: %v, raw response: %v, status code: %v", err, resp.RawResponse(), resp.Code)
	}

	symbols := common.TradableSymbols(resp.Data)

	pair_a := fmt.Sprintf("%s%s", t.cfg.Triangular.ABase, t.cfg.Triangular.AQuote)
	a_base := t.cfg.Triangular.ABase
	a_quote := t.cfg.Triangular.AQuote

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
