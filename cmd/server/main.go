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

func initLogger() (*zap.Logger, error){
    log, logErr := logger.NewLogger()
    if logErr != nil {
        return nil, logErr
    }
    sugar = log.Sugar()

    return log, nil
}

func getConf() {
    err := env.Parse(&envVars)

    if err != nil {
        sugar.Error("Couldn't get env vars")
        envVars.Address = "localhost:8080"
    }

    var fAddress = flag.String("a", envVars.Address, "enter IP format ip:port")

    flag.Parse()

    if envVars.Address == "localhost:8080" {
        envVars.Address = *fAddress
    }
}


func main() {
    log, err := initLogger()
    if err != nil {
        panic(err)
    }
    defer log.Sync()
    getConf()

    stor := storage.NewEmptyStorage()

    r := CreateRouter(stor)

    err = http.ListenAndServe(envVars.Address, r)

    if err != nil {
        sugar.Fatal(err)
    }
}
