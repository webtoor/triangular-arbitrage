package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/webtoor/triangular-arbitrage/internal/appctx"
	"github.com/webtoor/triangular-arbitrage/internal/arbitrage"
	arbitragebinance "github.com/webtoor/triangular-arbitrage/internal/arbitrage/binance"
	"github.com/webtoor/triangular-arbitrage/internal/bootstrap"
	"github.com/webtoor/triangular-arbitrage/internal/consts"
	"github.com/webtoor/triangular-arbitrage/internal/providers"
	"github.com/webtoor/triangular-arbitrage/internal/providers/binance"
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

	arb := arbitrage.New()
	arb.Registry(consts.Binance, arbitragebinance.New(cfg, spot))
	ctx := context.Background()

	for {
		select {
		case <-time.After(2 * time.Second):
			arbitrageTriangular := triangular.New(cfg, spot, arb)
			err := arbitrageTriangular.Start(ctx)
			if err != nil {
				logger.Error(err)
			}
		case <-ctx.Done():
			fmt.Println("We're done here!")
			return
		}
	}

}
