package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/webtoor/triangular-arbitrage/internal/appctx"
	"github.com/webtoor/triangular-arbitrage/internal/triangular"
	"github.com/webtoor/triangular-arbitrage/pkg/logger"
)

func Start(cfg *appctx.Config, t triangular.Resolverer) {
	var (
		lf = logger.NewFields(
			logger.EventName("scheduler.start"),
		)
	)

	logger.Info("starting triangular arbitrage", lf...)
	s := gocron.NewScheduler(time.Local)

	s.Every(cfg.Triangular.Interval).Second().Do(func() {
		err := t.Start(context.Background())
		if err != nil {
			logger.Error(fmt.Sprintf("scheduler stopped cause: %v", err), lf...)
			s.Stop()
		}
	})

	s.StartBlocking()
}
