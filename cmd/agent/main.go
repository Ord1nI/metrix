package main

import(
    "github.com/go-resty/resty/v2"
    "github.com/caarlos0/env/v11"

    "flag"
)

type Config struct {
    Address string `env:"ADDRESS"`
    ReportInterval int64 `env:"REPORT_INTERVAL"`
    PollInterval int64 `env:"POLL_INTERVAL"`
}

var metrics map[string]float64

var(
    envVars Config

    fIPStr *string
    fReportInterval *int64
    fPollInterval *int64
)

func init() {
    env.Parse(&envVars)

    if envVars.Address == "" {
        fIPStr = flag.String("a", "localhost:8080", "enter IP format ip:port")
    } else {
        fIPStr = &envVars.Address
    }

    if envVars.PollInterval == 0 {
        fPollInterval = flag.Int64("p", 2, "enter POOL INTERVAL in seconds")
    } else {
        fPollInterval = &envVars.PollInterval
    }

    if envVars.ReportInterval == 0 {
        fReportInterval = flag.Int64("r", 10, "enter REPORT INTERVAL in seconds")
    } else {
        fReportInterval = &envVars.ReportInterval
    }
}

func main() {
    flag.Parse()

    client := resty.New().SetBaseURL("http://"+*fIPStr)
    StartClient(client)

}
