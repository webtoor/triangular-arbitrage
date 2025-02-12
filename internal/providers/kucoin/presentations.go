package kucoin

type AllSymbolsResponse struct {
	Code string                   `json:"code"`
	Data []DataAllSymbolsResponse `json:"data"`
}

type DataAllSymbolsResponse struct {
	Symbol        string `json:"symbol"`
	BaseCurrency  string `json:"baseCurrency"`
	QuoteCurrency string `json:"quoteCurrency"`
	EnableTrading bool   `json:"enableTrading"`
}

type TickerPricesResponse struct {
	Code string                   `json:"code"`
	Data DataTickerPricesResponse `json:"data"`
}

type DataTickerPricesResponse struct {
	Time   int64               `json:"time"`
	Ticker []map[string]string `json:"ticker"`
}
