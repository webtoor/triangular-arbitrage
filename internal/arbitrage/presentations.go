package arbitrage

type PriceByTradingPairResp struct {
	PairA       string  `json:"pair_a"`
	PairB       string  `json:"pair_b"`
	PairC       string  `json:"pair_c"`
	PairAAsk    float64 `json:"pair_a_ask"`
	PairABid    float64 `json:"pair_a_bid"`
	PairBAsk    float64 `json:"pair_b_ask"`
	PairBBid    float64 `json:"pair_b_bid"`
	PairCAsk    float64 `json:"pair_c_ask"`
	PairCBid    float64 `json:"pair_c_bid"`
	PairAAskQty string  `json:"pair_a_ask_qty"`
	PairABidQty string  `json:"pair_a_bid_qty"`
	PairBAskQty string  `json:"pair_b_ask_qty"`
	PairBBidQty string  `json:"pair_b_bid_qty"`
	PairCAskQty string  `json:"pair_c_ask_qty"`
	PairCBidQty string  `json:"pair_c_bid_qty"`
}

type TriangularTradeParam struct {
	Exchange     string  `json:"exchange"`
	Direction    string  `json:"direction"`
	PairA        string  `json:"pair_a"`
	PairB        string  `json:"pair_b"`
	PairC        string  `json:"pair_c"`
	PricesA      float64 `json:"prices_a"`
	PricesB      float64 `json:"prices_b"`
	PricesC      float64 `json:"prices_c"`
	QtyPairA     float64 `json:"qty_pair_a"`
	QtyPairB     float64 `json:"qty_pair_b"`
	QtyPairC     float64 `json:"qty_pair_c"`
	InitialFunds float64 `json:"initial_funds"`
	FinalFunds   float64 `json:"final_funds"`
}
