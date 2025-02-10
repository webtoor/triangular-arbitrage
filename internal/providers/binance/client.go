package binance

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/webtoor/triangular-arbitrage/internal/appctx"
	"github.com/webtoor/triangular-arbitrage/internal/providers"
	"github.com/webtoor/triangular-arbitrage/pkg/httpx"
)

type binance struct {
	cfg *appctx.Config
}

func New(cfg *appctx.Config) providers.SpotAPI {
	return &binance{
		cfg: cfg,
	}
}

func (p *binance) ExchangeInfo(ctx context.Context) (appctx.Response, error) {
	var (
		resp     = appctx.NewResponse()
		respBody = AllSymbolsResponse{}
		url      = fmt.Sprintf("%s%s", p.cfg.Binance.BaseUrl, p.cfg.Binance.PathExchangeInfo)
	)

	reqOption := httpx.RequestOptions{
		Context: ctx,
		Method:  http.MethodGet,
		Timeout: time.Duration(p.cfg.Binance.Timeout) * time.Second,
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

func (p *binance) TickerPrices(ctx context.Context) (appctx.Response, error) {
	var (
		resp     = appctx.NewResponse()
		respBody = []providers.TickerPrices{}
		url      = fmt.Sprintf("%s%s", p.cfg.Binance.BaseUrl, p.cfg.Binance.PathTickerPrices)
	)

	reqOption := httpx.RequestOptions{
		Context: ctx,
		Method:  http.MethodGet,
		Timeout: time.Duration(p.cfg.Binance.Timeout) * time.Second,
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

func (p *binance) PlaceOrder(ctx context.Context, in any) (appctx.Response, error) {
	var (
		resp       = appctx.NewResponse()
		respBody   = map[string]any{}
		requestUrl = fmt.Sprintf("%s%s", p.cfg.Binance.BaseUrl, p.cfg.Binance.PathPlaceOrder)
	)

	param, ok := in.(PlaceOrderRequest)
	if !ok {
		return *resp.WithCode(http.StatusInternalServerError), fmt.Errorf("invalid parameter")
	}

	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	mapParams := url.Values{}
	mapParams.Set("symbol", param.Symbol)
	mapParams.Set("side", param.Side)
	mapParams.Set("type", fmt.Sprint(param.Type))
	mapParams.Set("quoteOrderQty", fmt.Sprint(param.Quantity))
	mapParams.Set("timestamp", fmt.Sprint(timestamp))
	signature := createSign(mapParams, p.cfg.Binance.SecretKey)
	mapParams.Set("signature", signature)

	requestUrl += "?"
	requestUrl += mapParams.Encode()

	h := httpx.Headers{}
	h.Add(httpx.XMBXAPIKEY, p.cfg.Binance.ApiKey)

	reqOption := httpx.RequestOptions{
		Context: ctx,
		Method:  http.MethodPost,
		Timeout: time.Duration(p.cfg.Binance.Timeout) * time.Second,
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

	return *resp.WithCode(req.Status()).WithData(respBody), nil
}
