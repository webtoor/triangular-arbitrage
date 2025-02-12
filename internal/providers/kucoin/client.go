package kucoin

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/webtoor/triangular-arbitrage/internal/appctx"
	"github.com/webtoor/triangular-arbitrage/internal/providers"
	"github.com/webtoor/triangular-arbitrage/pkg/httpx"
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
	return *appctx.NewResponse(), fmt.Errorf("invalid exchange")
}

func (p *kucoin) PlaceOrder(ctx context.Context, in any) (appctx.Response, error) {
	return *appctx.NewResponse(), fmt.Errorf("invalid exchange")
}
