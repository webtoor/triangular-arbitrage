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
	ID          string  `json:"id,omitempty"`
	Symbol      string  `json:"symbol,omitempty"`
	Side        string  `json:"side,omitempty"`
	Type        string  `json:"type,omitempty"`
	TimeInForce string  `json:"timeInForce,omitempty"`
	Quantity    float64 `json:"quantity,omitempty"`
}
