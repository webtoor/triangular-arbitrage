package binance

import (
	"context"
	"errors"
	"fmt"
	"net/http"
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
		return *resp.WithCode(500), fmt.Errorf("error: %v", err.Error())
	}

	if req.Status() != http.StatusOK {
		return *resp.WithCode(req.Status()).WithRawResponse(req.String()), fmt.Errorf("http status code %v", req.Status())
	}

	err = req.DecodeJSON(&respBody)

	if err != nil {
		return *resp.WithCode(500), errors.New("error decode json")
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
		return *resp.WithCode(500), fmt.Errorf("error: %v", err.Error())
	}

	if req.Status() != http.StatusOK {
		return *resp.WithCode(req.Status()).WithRawResponse(req.String()), fmt.Errorf("http status code %v", req.Status())
	}

	err = req.DecodeJSON(&respBody)

	if err != nil {
		return *resp.WithCode(500), errors.New("error decode json")
	}

	return *resp.WithCode(req.Status()).WithData(respBody), nil

}
