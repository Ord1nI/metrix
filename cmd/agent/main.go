package main

import(
    "github.com/go-resty/resty/v2"
    "flag"
)

var metrics map[string]float64

var(
    fIpStr = flag.String("a", "http://localhost:8080", "enter IP format ip:port")
    fReportInterval = flag.Int64("r", 10, "enter REPORT INTERVAL in seconds")
    fPollInterval = flag.Int64("p", 2, "enter POOL INTERVAL in seconds")
)

func main() {
    flag.Parse()

    client := resty.New().SetBaseURL(*fIpStr)
    StartClient(client)

}
