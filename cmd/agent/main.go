package main

import (
	"github.com/go-resty/resty/v2"

	"github.com/Ord1nI/metrix/internal/logger"
	"github.com/Ord1nI/metrix/internal/repo/storage"
	"github.com/Ord1nI/metrix/internal/configs"
)


var (
    envVars = configs.AgentConfig{}
    sugar logger.Logger
)

func initF() {
    log, err := logger.New()
    if err != nil {
        panic(err)
    }
    sugar = log
    log.Infoln("loger created successfuly")

    configs.GetAgentConf(sugar, &envVars)
    log.Infoln("sysvars:", envVars)
}



func main() {
    initF()

    stor := storage.NewMemStorage()

    client := resty.New().SetBaseURL("http://" + envVars.Address)
    StartClient(client, stor)
}
