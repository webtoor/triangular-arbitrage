package scheduler

import (
	"context"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/webtoor/triangular-arbitrage/internal/appctx"
	"github.com/webtoor/triangular-arbitrage/internal/triangular"
	"github.com/webtoor/triangular-arbitrage/pkg/logger"
)

func Start(cfg *appctx.Config, t triangular.Resolve) {
	s := gocron.NewScheduler(time.Local)

	s.Every(cfg.Triangular.Interval).Second().Do(func() {
		if cfg.Binance.TriangularEnabled {
			err := t.Start(context.Background())
			if err != nil {
				logger.Error(err)
			}
		}
	})

	s.StartBlocking()
}
