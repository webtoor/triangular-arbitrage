package arbitrage

type PriceByTradingPairResp struct {
	PairAAsk  float64 `json:"pair_a_ask"`
	PairABid  float64 `json:"pair_a_bid"`
	PairBBid  float64 `json:"pair_b_bid"`
	PairBAAsk float64 `json:"pair_b_ask"`
	PairCAAsk float64 `json:"pair_c_ask"`
	PairCBid  float64 `json:"pair_c_bid"`
}
