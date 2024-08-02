package main

import(
    "github.com/go-resty/resty/v2"
    "github.com/caarlos0/env/v11"

    "github.com/Ord1nI/metrix/internal/storage"
    "flag"
)

type Config struct {
    Address string `env:"ADDRESS" envDefault:"localhost:8080"` //envvar $ADDRESS or envDefault
    ReportInterval int64 `env:"REPORT_INTERVAL" envDefault:"10"` //envvar $REPORTINTERVAL or envDefault
    PollInterval int64 `env:"POLL_INTERVAL" envDefault:"2"`    //envvar $POOLINTERVAL or envDefault
}

var (
    envVars Config
)

func init() {
    err := env.Parse(&envVars)

    if err != nil {
        panic(err)
    }

    if envVars.Address ==  "localhost:8080" {
        envVars.Address = *flag.String("a", envVars.Address, "enter IP format ip:port")
    }

    if envVars.PollInterval == 2 {
        envVars.PollInterval = *flag.Int64("p", envVars.PollInterval, "enter POOL INTERVAL in seconds")
    }

    if envVars.ReportInterval == 10 {
        envVars.ReportInterval = *flag.Int64("r", envVars.ReportInterval, "enter REPORT INTERVAL in seconds")
    }
}

func main() {
    flag.Parse()
    
    stor := storage.NewEmptyStorage()

    client := resty.New().SetBaseURL("http://" + envVars.Address)
    StartClient(client, stor)
}
