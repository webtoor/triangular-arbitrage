package kucoin

type AllSymbolsResponse struct {
	Code string `json:"code"`
	Data []Data `json:"data"`
}

type Data struct {
	Symbol        string `json:"symbol"`
	BaseCurrency  string `json:"baseCurrency"`
	QuoteCurrency string `json:"quoteCurrency"`
	EnableTrading bool   `json:"enableTrading"`
}
