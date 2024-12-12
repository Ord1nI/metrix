
package agent

import (
	"encoding/json"
	"flag"
	"io"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Address         string `json:"address" env:"ADDRESS" envDefault:"localhost:8080"`
	Key             string `json:"key" env:"KEY" envDefault:""`
	PublicKeyFile   string `json:"crypto_key" env:"CRYPTO_KEY" envDefault:""`
	FileCfg         string `env:"AGENT_CONFIG" envDefault:""`
	BackoffSchedule []time.Duration
	PollInterval    int64 `json:"poll_interval" env:"POLL_INTERVAL" envDefault:"2"`
	ReportInterval  int64 `json:"report_interval" env:"REPORT_INTERVAL" envDefault:"10"`
	RateLimit       int64   `json:"rate_limit" env:"RATE_LIMIT" envDefault:"1"`
}

func setValue[i int64 | string] (mainCfg *i, fileCfg *i, flag *i, defaultValue i) {
	if *mainCfg == defaultValue {
		*mainCfg = *flag;
	}

	var sValue i;

	if *mainCfg == defaultValue && *fileCfg != sValue{
		*mainCfg = *fileCfg
	}
}

func (a *Agent) getConfFromFile(fileName string) Config {

	file, err := os.OpenFile(fileName, os.O_RDONLY, 0644);

	if err != nil {
		a.Logger.Fatalln(err);
	}

	var cfg = Config{}

	bFile, err := io.ReadAll(file)

	if err != nil {
		a.Logger.Fatalln(err)
	}

	err = json.Unmarshal(bFile, &cfg)

	if err != nil {
		a.Logger.Fatalln(err)
	}

	return cfg
}

func (a *Agent) GetConf() {
	err := env.Parse(&a.Config)

	if err != nil {
		a.Logger.Fatalln(err)
	}

	var (
		fAddress = flag.String("a", a.Config.Address, "enter IP format ip:port")

		fPoolInterval = flag.Int64("p", a.Config.PollInterval, "enter POOL INTERVAL in seconds")

		fReportInterval = flag.Int64("r", a.Config.ReportInterval, "enter REPORT INTERVAL in seconds")

		fKey = flag.String("k", a.Config.Key, "enter Signatur key")

		fFileCfg = flag.String("config", a.Config.FileCfg, "enter config file location")

		fPublicKeyFile = flag.String("crypto-key", a.Config.PublicKeyFile, "enter location of file with public key")

		fRateLimit = flag.Int64("l", a.Config.RateLimit, "enter Rate limit")
	)

	flag.Parse()

	if a.Config.FileCfg == "" {
		a.Config.FileCfg = *fFileCfg
	}

	var cfgFromFile Config

	if a.Config.FileCfg != ""  {
		cfgFromFile = a.getConfFromFile(a.Config.FileCfg)
	}


	setValue (&a.Config.Address, &cfgFromFile.Address, fAddress, "localhost:8080")
	setValue (&a.Config.PollInterval, &cfgFromFile.PollInterval, fPoolInterval, 2)
	setValue (&a.Config.ReportInterval, &cfgFromFile.ReportInterval, fReportInterval, 10)
	setValue (&a.Config.Key, &cfgFromFile.Key, fKey, "")
	setValue (&a.Config.PublicKeyFile, &cfgFromFile.PublicKeyFile, fPublicKeyFile, "")
	setValue (&a.Config.RateLimit, &cfgFromFile.RateLimit, fRateLimit, 1)

	a.Config.BackoffSchedule = []time.Duration{
		100 * time.Millisecond,
		500 * time.Millisecond,
		1 * time.Second,
	}

	if a.Config.PollInterval > a.Config.ReportInterval {
		a.Logger.Infoln("PollInterval > ReportInterval")
	}
}
