package scheduler

import (
	"context"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/webtoor/triangular-arbitrage/internal/appctx"
	"github.com/webtoor/triangular-arbitrage/internal/triangular"
)

func Start(cfg *appctx.Config, t triangular.Resolverer) {
	s := gocron.NewScheduler(time.Local)

	s.Every(cfg.Triangular.Interval).Second().Do(func() {
		t.Start(context.Background())
	})

	s.StartBlocking()
}
