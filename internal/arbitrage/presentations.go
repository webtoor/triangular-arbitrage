package arbitrage

type PriceByTradingPairResp struct {
	PairAAsk float64 `json:"pair_a_ask"`
	PairABid float64 `json:"pair_a_bid"`
	PairBBid float64 `json:"pair_b_bid"`
	PairBAsk float64 `json:"pair_b_ask"`
	PairCAsk float64 `json:"pair_c_ask"`
	PairCBid float64 `json:"pair_c_bid"`
}

type TriangularTradeParam struct {
	Direction    string  `json:"direction"`
	PairA        string  `json:"pair_a"`
	PairB        string  `json:"pair_b"`
	PairC        string  `json:"pair_c"`
	FinalBalance float64 `json:"final_balance"`
}
