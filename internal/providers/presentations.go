package providers

type TickerPrices struct {
	Symbol    string `json:"symbol"`
	LastPrice string `json:"lastPrice"`
	BidPrice  string `json:"bidPrice"`
	BidQty    string `json:"bidQty"`
	AskPrice  string `json:"askPrice"`
	AskQty    string `json:"askQty"`
}

type PlaceOrderRequest struct {
	Symbol      string  `json:"symbol"`
	Side        string  `json:"side"`
	Type        string  `json:"type"`
	TimeInForce string  `json:"timeInForce"`
	Quantity    float64 `json:"quantity"`
}
