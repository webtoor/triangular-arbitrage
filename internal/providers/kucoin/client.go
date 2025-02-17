package kucoin

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/webtoor/triangular-arbitrage/internal/appctx"
	"github.com/webtoor/triangular-arbitrage/internal/providers"
	"github.com/webtoor/triangular-arbitrage/pkg/httpx"
	"github.com/webtoor/triangular-arbitrage/pkg/util"
)

type kucoin struct {
	cfg *appctx.Config
}

func New(cfg *appctx.Config) providers.SpotAPI {
	return &kucoin{
		cfg: cfg,
	}
}

func (p *kucoin) GetAllSymbols(ctx context.Context) (appctx.Response, error) {
	var (
		resp     = appctx.NewResponse()
		respBody = AllSymbolsResponse{}
		url      = fmt.Sprintf("%s%s", p.cfg.Kucoin.BaseUrl, p.cfg.Kucoin.PathSymbols)
	)

	reqOption := httpx.RequestOptions{
		Context: ctx,
		Method:  http.MethodGet,
		Timeout: time.Duration(p.cfg.Kucoin.Timeout) * time.Second,
		URL:     url,
	}

	req, err := httpx.Request(reqOption)

	if err != nil {
		return *resp.WithCode(http.StatusInternalServerError), fmt.Errorf("error: %v", err)
	}

	if req.Status() != http.StatusOK {
		return *resp.WithCode(req.Status()).WithRawResponse(req.String()), fmt.Errorf("http status code %v", req.Status())
	}

	err = req.DecodeJSON(&respBody)

	if err != nil {
		return *resp.WithCode(http.StatusInternalServerError), fmt.Errorf("error: %v", err)
	}

	return *resp.WithCode(req.Status()).WithData(respBody), nil
}

func (p *kucoin) TickerPrices(ctx context.Context) (appctx.Response, error) {
	var (
		resp     = appctx.NewResponse()
		tickers  []providers.TickerPrices
		respBody TickerPricesResponse
		url      = fmt.Sprintf("%s%s", p.cfg.Kucoin.BaseUrl, p.cfg.Kucoin.PathTickerPrices)
	)

	reqOption := httpx.RequestOptions{
		Context: ctx,
		Method:  http.MethodGet,
		Timeout: time.Duration(p.cfg.Kucoin.Timeout) * time.Second,
		URL:     url,
	}

	req, err := httpx.Request(reqOption)

	if err != nil {
		return *resp.WithCode(http.StatusInternalServerError), fmt.Errorf("error: %v", err)
	}

	if req.Status() != http.StatusOK {
		return *resp.WithCode(req.Status()).WithRawResponse(req.String()), fmt.Errorf("http status code %v", req.Status())
	}

	err = req.DecodeJSON(&respBody)

	if err != nil {
		return *resp.WithCode(http.StatusInternalServerError), fmt.Errorf("error: %v", err)
	}

	for _, v := range respBody.Data.Ticker {
		tickers = append(tickers, providers.TickerPrices{
			Symbol:    v["symbol"],
			LastPrice: v["last"],
			BidPrice:  v["buy"],
			BidQty:    v["bestBidSize"],
			AskPrice:  v["sell"],
			AskQty:    v["bestAskSize"],
		})
	}
	return *resp.WithCode(req.Status()).WithData(tickers), nil
}

func (p *kucoin) GetOrderBook(ctx context.Context, symbol string) (appctx.Response, error) {
	return *appctx.NewResponse(), fmt.Errorf("GetOrderBook: invalid exchange")
}

func (p *kucoin) PlaceOrder(ctx context.Context, in any) (appctx.Response, error) {
	var (
		resp       = appctx.NewResponse()
		respBody   map[string]any
		requestUrl = fmt.Sprintf("%s%s", p.cfg.Kucoin.BaseUrl, p.cfg.Kucoin.PathPlaceOrder)
	)

	param, ok := in.(providers.PlaceOrderRequest)
	if !ok {
		return *resp.WithCode(http.StatusInternalServerError), fmt.Errorf("invalid parameter")
	}

	pl := PlaceOrderRequest{
		Symbol: param.Symbol,
		Type:   strings.ToLower(param.Type),
		Side:   strings.ToLower(param.Side),
		Size:   cast.ToString(param.Quantity),
	}

	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	sign := createSign(fmt.Sprintf("%s%s%s", http.MethodPost, p.cfg.Kucoin.PathPlaceOrder, util.ToJSON(pl)), cast.ToString(timestamp), p.cfg.Kucoin.SecretKey)
	passphrase := createSign(p.cfg.Kucoin.ApiPassphrase, util.EmptyString(), p.cfg.Kucoin.SecretKey)

	h := httpx.Headers{}
	h.Add("KC-API-KEY", p.cfg.Kucoin.ApiKey)
	h.Add("KC-API-PASSPHRASE", passphrase)
	h.Add("KC-API-TIMESTAMP", cast.ToString(timestamp))
	h.Add("KC-API-SIGN", sign)
	h.Add("KC-API-KEY-VERSION", "2")
	h.Add("Content-Type", "application/json")

	reqOption := httpx.RequestOptions{
		Context: ctx,
		Method:  http.MethodPost,
		Timeout: time.Duration(p.cfg.Kucoin.Timeout) * time.Second,
		URL:     requestUrl,
		Payload: pl,
		Header:  h,
	}

	req, err := httpx.Request(reqOption)

	if err != nil {
		return *resp.WithCode(http.StatusInternalServerError), fmt.Errorf("%v", err)
	}

	if req.Status() != http.StatusOK {
		return *resp.WithCode(req.Status()).WithRawResponse(req.String()), fmt.Errorf("http status code %v", req.Status())
	}

	err = req.DecodeJSON(&respBody)

	if err != nil {
		return *resp.WithCode(http.StatusInternalServerError), fmt.Errorf("%v", err)
	}

	_, ok = respBody["data"]

	if !ok {
		return *resp.WithCode(http.StatusInternalServerError).WithRawResponse(req.String()), fmt.Errorf("%v", respBody["msg"])
	}

	return *resp.WithCode(req.Status()).WithData(respBody).WithRawResponse(req.String()), nil
}

func (p *kucoin) GetOrderByID(ctx context.Context, id, symbol string) (appctx.Response, error) {
	var (
		resp       = appctx.NewResponse()
		respBody   = GetOrderByIDResponse{}
		requestUrl = fmt.Sprintf("%s%s/%s", p.cfg.Kucoin.BaseUrl, p.cfg.Kucoin.PathGetOrderByID, id)
	)

	mapParams := url.Values{}
	mapParams.Set("symbol", symbol)

	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	sign := createSign(fmt.Sprintf("%s%s/%s?%s", http.MethodGet, p.cfg.Kucoin.PathGetOrderByID, id, mapParams.Encode()), cast.ToString(timestamp), p.cfg.Kucoin.SecretKey)
	passphrase := createSign(p.cfg.Kucoin.ApiPassphrase, util.EmptyString(), p.cfg.Kucoin.SecretKey)

	h := httpx.Headers{}
	h.Add("KC-API-KEY", p.cfg.Kucoin.ApiKey)
	h.Add("KC-API-PASSPHRASE", passphrase)
	h.Add("KC-API-TIMESTAMP", cast.ToString(timestamp))
	h.Add("KC-API-SIGN", sign)
	h.Add("KC-API-KEY-VERSION", "3")

	requestUrl += "?"
	requestUrl += mapParams.Encode()

	reqOption := httpx.RequestOptions{
		Context: ctx,
		Method:  http.MethodGet,
		Timeout: time.Duration(p.cfg.Kucoin.Timeout) * time.Second,
		URL:     requestUrl,
		Header:  h,
	}

	req, err := httpx.Request(reqOption)

	if err != nil {
		return *resp.WithCode(http.StatusInternalServerError), fmt.Errorf("error: %v", err)
	}

	if req.Status() != http.StatusOK {
		return *resp.WithCode(req.Status()).WithRawResponse(req.String()), fmt.Errorf("http status code %v", req.Status())
	}

	err = req.DecodeJSON(&respBody)

	if err != nil {
		return *resp.WithCode(http.StatusInternalServerError), fmt.Errorf("error: %v", err)
	}

	return *resp.WithCode(req.Status()).WithData(respBody).WithRawResponse(req.String()), nil
}
