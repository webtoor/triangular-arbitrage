package cmd

import (
	"github.com/webtoor/triangular-arbitrage/cmd/scheduler"
	"github.com/webtoor/triangular-arbitrage/internal/appctx"
	"github.com/webtoor/triangular-arbitrage/internal/arbitrage"
	arbitragebinance "github.com/webtoor/triangular-arbitrage/internal/arbitrage/binance"
	arbitragekucoin "github.com/webtoor/triangular-arbitrage/internal/arbitrage/kucoin"
	"github.com/webtoor/triangular-arbitrage/internal/bootstrap"
	"github.com/webtoor/triangular-arbitrage/internal/consts"
	"github.com/webtoor/triangular-arbitrage/internal/providers"
	"github.com/webtoor/triangular-arbitrage/internal/providers/binance"
	"github.com/webtoor/triangular-arbitrage/internal/providers/kucoin"
	"github.com/webtoor/triangular-arbitrage/internal/triangular"
	"github.com/webtoor/triangular-arbitrage/pkg/logger"
)

func Start() {
	cfg := appctx.NewConfig()
	bootstrap.RegistryLogger(cfg)
	logger.SetJSONFormatter()

	// init providers
	spot := providers.New()
	spot.Registry(consts.Binance, binance.New(cfg))
	spot.Registry(consts.Kucoin, kucoin.New(cfg))

	arb := arbitrage.New()
	arb.Registry(consts.Binance, arbitragebinance.New(cfg, spot))
	arb.Registry(consts.Kucoin, arbitragekucoin.New(cfg, spot))
	arbitrageTriangular := triangular.New(cfg, spot, arb)

	scheduler.Start(cfg, arbitrageTriangular)
}
