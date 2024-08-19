package main

import (
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"


	"github.com/Ord1nI/metrix/internal/logger"
	"github.com/Ord1nI/metrix/internal/repo/storage"
	"github.com/Ord1nI/metrix/internal/configs"
)


var (
    envVars = configs.AgentConfig{}
    sugar *zap.SugaredLogger
)

func initF() {
    log, err := logger.NewLogger()
    if err != nil {
        panic(err)
    }
    defer log.Sync()
    sugar = log.Sugar()
    sugar.Infoln("loger created successfuly")

    configs.GetAgentConf(sugar, &envVars)
}



func main() {
    initF()

    stor := storage.NewMemStorage()

    client := resty.New().SetBaseURL("http://" + envVars.Address)
    StartClient(client, stor)
}
