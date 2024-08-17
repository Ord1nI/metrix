package configs

import(
    "time"
	"flag"

	"github.com/Ord1nI/metrix/internal/logger"
	"github.com/caarlos0/env/v11"
)

type AgentConfig struct {
    Address string `env:"ADDRESS" envDefault:"localhost:8080"` //envvar $ADDRESS or envDefault
    PollInterval int64 `env:"POLL_INTERVAL" envDefault:"2"`    //envvar $POOLINTERVAL or envDefault
    ReportInterval int64 `env:"REPORT_INTERVAL" envDefault:"10"` //envvar $REPORTINTERVAL or envDefault
    BackoffSchedule []time.Duration
}


func GetAgentConf(sugar logger.Logger, envVars *AgentConfig) {
    err := env.Parse(envVars)

    if err != nil {
        sugar.Errorln(err)

        envVars.Address = "localhost:8080"
        envVars.PollInterval = 2
        envVars.ReportInterval = 10
    }

    var (
        fAddress = flag.String("a", envVars.Address, "enter IP format ip:port")

        fPoolInterval = flag.Int64("p", envVars.PollInterval, "enter POOL INTERVAL in seconds")

        fReportInterval = flag.Int64("r", envVars.ReportInterval, "enter REPORT INTERVAL in seconds")
    )

    flag.Parse()

    if envVars.Address == "localhost:8080" {
        envVars.Address = *fAddress
    }

    if envVars.PollInterval == 2 {
        envVars.PollInterval = *fPoolInterval
    }

    if envVars.ReportInterval == 10 {
        envVars.ReportInterval = *fReportInterval
    }
}
