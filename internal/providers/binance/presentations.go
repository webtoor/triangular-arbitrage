package binance

type AllSymbolsResponse struct {
	TimeZone   string    `json:"timeZone"`
	ServerTime int64     `json:"serverTime"`
	Symbols    []Symbols `json:"symbols"`
}

type Symbols struct {
	Symbol               string `json:"symbol"`
	Status               string `json:"status"`
	BaseAsset            string `json:"baseAsset"`
	QuoteAsset           string `json:"quoteAsset"`
	IsSpotTradingAllowed bool   `json:"isSpotTradingAllowed"`
}

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
