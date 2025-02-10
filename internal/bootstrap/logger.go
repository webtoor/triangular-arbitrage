package bootstrap

import (
	"github.com/webtoor/triangular-arbitrage/internal/appctx"
	"github.com/webtoor/triangular-arbitrage/pkg/logger"
	"github.com/webtoor/triangular-arbitrage/pkg/util"
)

func RegistryLogger(cfg *appctx.Config) {
	logger.Setup(logger.Config{
		Environment: util.EnvironmentTransform(cfg.App.Env),
		Debug:       cfg.App.Debug,
		Level:       cfg.Logger.Level,
		ServiceName: cfg.Logger.Name,
	})
}
