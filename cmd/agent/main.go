package main

import(
    "github.com/go-resty/resty/v2"
    "github.com/caarlos0/env/v11"
    "go.uber.org/zap"

    "github.com/Ord1nI/metrix/internal/storage"
    "github.com/Ord1nI/metrix/internal/logger"
    "flag"
)

type Config struct {
    Address string `env:"ADDRESS" envDefault:"localhost:8080"` //envvar $ADDRESS or envDefault
    PollInterval int64 `env:"POLL_INTERVAL" envDefault:"2"`    //envvar $POOLINTERVAL or envDefault
    ReportInterval int64 `env:"REPORT_INTERVAL" envDefault:"10"` //envvar $REPORTINTERVAL or envDefault
}

var (
    envVars Config
    sugar *zap.SugaredLogger
)

func initLogger(log *zap.Logger) error{
    log, logErr := logger.NewLogger()
    if logErr != nil {
        return logErr
    }
    sugar = log.Sugar()

    return nil
}

func getConf() {
    err := env.Parse(&envVars)

    if err != nil {
        sugar.Error(err)

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


func main() {
    var log *zap.Logger

    initLogger(log)
    
    defer log.Sync()

    getConf()

    
    stor := storage.NewMemStorage()

    client := resty.New().SetBaseURL("http://" + envVars.Address)
    StartClient(client, stor)
}
