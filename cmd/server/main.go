package main

import (
    "github.com/caarlos0/env/v11"
    "go.uber.org/zap"

    "flag"
    "net/http"
    "github.com/Ord1nI/metrix/internal/storage"
    "github.com/Ord1nI/metrix/internal/logger"
)

type Config struct {
    Address string `env:"ADDRESS" envDefault:"localhost:8080"` //envvar $ADDRESS or envDefault
}

var envVars Config

var sugar *zap.SugaredLogger

func getConf() {
    err := env.Parse(&envVars)

    if err != nil {
        panic(err)
    }

    var fAddress = flag.String("a", envVars.Address, "enter IP format ip:port")

    flag.Parse()

    if envVars.Address == "localhost:8080" {
        envVars.Address = *fAddress
    }
}


func main() {
    getConf()

    logger, logErr := logger.NewLogger()
    logger.WithOptions(zap.AddCaller())
    if logErr != nil {
        panic(logErr)
    }

    defer logger.Sync()

    sugar = logger.Sugar()

    stor := storage.NewEmptyStorage()

    r := CreateRouter(stor)

    err := http.ListenAndServe(envVars.Address, r)

    if err != nil {
        panic(err)
    }
}
