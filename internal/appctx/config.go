package appctx

import (
	"fmt"
	"log"
	"sync"

	"github.com/webtoor/triangular-arbitrage/internal/consts"
	"github.com/webtoor/triangular-arbitrage/pkg/file"
)

var (
	once sync.Once
	_cfg *Config
)

type Config struct {
	App             *Common    `yaml:"app" json:"app"`
	Triangular      Triangular `yaml:"triangular" json:"triangular"`
	TriangularPairs []TriangularBinancePair
	Logger          Logging `yaml:"logger" json:"logger"`
	Binance         Binance `yaml:"binance" json:"binance"`
	Kucoin          Kucoin  `yaml:"kucoin" json:"kucoin"`
}

type Common struct {
	AppName  string `yaml:"name" json:"name"`
	Debug    bool   `yaml:"debug" json:"debug"`
	Timezone string `yaml:"timezone" json:"timezone"`
	Env      string `yaml:"env" json:"env"`
}

type Triangular struct {
	Exchanges []string `yaml:"exchanges" json:"exchanges"`
	Interval  int      `yaml:"interval" json:"interval"`
	Balance   float64  `yaml:"balance" json:"balance"`
}

type Logging struct {
	Name  string `yaml:"name" json:"name"`
	Level string `yaml:"level" json:"level"`
}

type Binance struct {
	ApiKey              string  `yaml:"api_key" json:"api_key"`
	SecretKey           string  `yaml:"secret_key" json:"secret_key"`
	PathTriangularPairs string  `yaml:"path_triangular_pairs" json:"path_triangular_pairs"`
	TriangularEnabled   bool    `yaml:"triangular_enabled" json:"triangular_enabled"`
	TradeEnabled        bool    `yaml:"trade_enabled" json:"trade_enabled"`
	Fees                float64 `yaml:"fees" json:"fees"`
	BaseUrl             string  `yaml:"base_url" json:"base_url"`
	PathExchangeInfo    string  `yaml:"path_exchange_info" json:"path_exchange_info"`
	PathTickerPrices    string  `yaml:"path_ticker_prices" json:"path_ticker_prices"`
	PathPlaceOrder      string  `yaml:"path_place_order" json:"path_place_order"`
	Timeout             int     `yaml:"timeout" json:"timeout"`
}

type Kucoin struct {
	ApiKey              string  `yaml:"api_key" json:"api_key"`
	SecretKey           string  `yaml:"secret_key" json:"secret_key"`
	PathTriangularPairs string  `yaml:"path_triangular_pairs" json:"path_triangular_pairs"`
	TriangularEnabled   bool    `yaml:"triangular_enabled" json:"triangular_enabled"`
	TradeEnabled        bool    `yaml:"trade_enabled" json:"trade_enabled"`
	Fees                float64 `yaml:"fees" json:"fees"`
	BaseUrl             string  `yaml:"base_url" json:"base_url"`
	PathSymbols         string  `yaml:"path_symbols" json:"path_symbols"`
	PathTickerPrices    string  `yaml:"path_ticker_prices" json:"path_ticker_prices"`
	PathPlaceOrder      string  `yaml:"path_place_order" json:"path_place_order"`
	Timeout             int     `yaml:"timeout" json:"timeout"`
}

type TriangularBinancePair struct {
	ABase    string `json:"a_base"`
	AQuote   string `json:"a_quote"`
	BBase    string `json:"b_base"`
	BQuote   string `json:"b_quote"`
	CBase    string `json:"c_base"`
	CQuote   string `json:"c_quote"`
	PairA    string `json:"pair_a"`
	PairB    string `json:"pair_b"`
	PairC    string `json:"pair_c"`
	Combined string `json:"combined"`
}

type TriangularKucoinPair struct {
	ABase    string `json:"a_base"`
	AQuote   string `json:"a_quote"`
	BBase    string `json:"b_base"`
	BQuote   string `json:"b_quote"`
	CBase    string `json:"c_base"`
	CQuote   string `json:"c_quote"`
	PairA    string `json:"pair_a"`
	PairB    string `json:"pair_b"`
	PairC    string `json:"pair_c"`
	Combined string `json:"combined"`
}

func NewConfig() *Config {
	fpath := []string{consts.ConfigPath}
	once.Do(func() {
		c, err := readCfg("app.yaml", fpath...)
		if err != nil {
			log.Fatal(err)
		}

		tp, err := readPairs("binance-pairs.json", fpath...)
		if err != nil {
			log.Fatal(err)
		}

		c.TriangularPairs = tp
		_cfg = c
	})

	return _cfg
}

func readCfg(fname string, ps ...string) (*Config, error) {
	var cfg *Config
	var errs []error

	for _, p := range ps {
		f := fmt.Sprint(p, fname)
		err := file.ReadFromYAML(f, &cfg)
		if err != nil {
			errs = append(errs, fmt.Errorf("file %s error %s", f, err.Error()))
			continue
		}
		break
	}

	if cfg == nil {
		return nil, fmt.Errorf("file config parse error %v", errs)
	}

	return cfg, nil
}

func readPairs(fname string, ps ...string) ([]TriangularBinancePair, error) {
	var tp []TriangularBinancePair
	var errs []error

	for _, p := range ps {
		f := fmt.Sprint(p, fname)
		err := file.ReadFromJSON(f, &tp)
		if err != nil {
			errs = append(errs, fmt.Errorf("file %s error %s", f, err.Error()))
			continue
		}
		break
	}

	if tp == nil {
		return nil, fmt.Errorf("file triangular pairs parse error %v", errs)
	}

	return tp, nil
}
