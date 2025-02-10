package providers

type TickerPrices struct {
	Symbol    string `json:"symbol"`
	LastPrice string `json:"lastPrice"`
	BidPrice  string `json:"bidPrice"`
	BidQty    string `json:"bidQty"`
	AskPrice  string `json:"askPrice"`
	AskQty    string `json:"askQty"`
}
