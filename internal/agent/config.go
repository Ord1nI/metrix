package agent

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Address         string `env:"ADDRESS" envDefault:"localhost:8080"` //envvar $ADDRESS or envDefault
	PollInterval    int64  `env:"POLL_INTERVAL" envDefault:"2"`        //envvar $POOLINTERVAL or envDefault
	ReportInterval  int64  `env:"REPORT_INTERVAL" envDefault:"10"`     //envvar $REPORTINTERVAL or envDefault
	BackoffSchedule []time.Duration
}

func (a *Agent) GetConf() {
	err := env.Parse(&a.Config)

	if err != nil {
		a.Logger.Errorln(err)

		a.Config.Address = "localhost:8080"
		a.Config.PollInterval = 2
		a.Config.ReportInterval = 10
	}

	var (
		fAddress = flag.String("a", a.Config.Address, "enter IP format ip:port")

		fPoolInterval = flag.Int64("p", a.Config.PollInterval, "enter POOL INTERVAL in seconds")

		fReportInterval = flag.Int64("r", a.Config.ReportInterval, "enter REPORT INTERVAL in seconds")
	)

	flag.Parse()

	if a.Config.Address == "localhost:8080" {
		a.Config.Address = *fAddress
	}

	if a.Config.PollInterval == 2 {
		a.Config.PollInterval = *fPoolInterval
	}

	if a.Config.ReportInterval == 10 {
		a.Config.ReportInterval = *fReportInterval
	}
	a.Config.BackoffSchedule = []time.Duration{
		100 * time.Millisecond,
		500 * time.Millisecond,
		1 * time.Second,
	}
	if a.Config.PollInterval > a.Config.ReportInterval {
		a.Logger.Infoln("PollInterval > ReportInterval")
	}
}
