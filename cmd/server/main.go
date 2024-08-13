package main

import (
    "github.com/caarlos0/env/v11"
    "go.uber.org/zap"

    "flag"
    "time"
    "net/http"
    "github.com/Ord1nI/metrix/internal/storage"
    "github.com/Ord1nI/metrix/internal/logger"
)

type Config struct {
    Address string `env:"ADDRESS" envDefault:"localhost:8080"` //envvar $ADDRESS or envDefault
    StoreInterval int `env:"STORE_INTERVAL" envDefault:"300"`
    FileStoragePath string `env:"FILE_STORAGE_PATH"`
    Restore bool `env:"RESTORE" envDefault:"true"`
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
    var fStoreInterval = flag.Int("i", envVars.StoreInterval,
        "enter interval (in seconds) between all data saved to specified file")
    var fFileStoragePath = flag.String("f", envVars.FileStoragePath,
        "enter path to file where all data will be saved")
    var fRestore = flag.Bool("r", envVars.Restore, "whether or not load data to specified file")

    flag.Parse()

    if envVars.Address == "localhost:8080" {
        envVars.Address = *fAddress
    }
    if envVars.StoreInterval == 300 {
        envVars.StoreInterval = *fStoreInterval
    }
    if envVars.FileStoragePath == "" {
        envVars.FileStoragePath = *fFileStoragePath
    }
    if envVars.Restore == true {
        envVars.Restore = *fRestore
    }
}

func StartDataSaver(s *storage.MemStorage) {
    for true {
        time.Sleep(time.Duration(envVars.StoreInterval) * time.Second)
        err := s.WriteToFile(envVars.FileStoragePath)
        if err != nil {
            sugar.Fatal(err)
        } esle {
            sugar.Info("Data saved")
        }
    }
}

func main() {
    log, err := initLogger()
    if err != nil {
        panic(err)
    }
    defer log.Sync()
    getConf()

    stor := storage.NewMemStorage()

    if envVars.FileStoragePath != "" && envVars.Restore == true {
        err := stor.GetFromFile(envVars.FileStoragePath)
        if err != nil {
            sugar.Info(err)
        } else {
            sugar.Info("Data loaded succesful",stor.Gauge)
        }
    }

    r := CreateRouter(stor)

    if  envVars.FileStoragePath != "" && 
        envVars.StoreInterval != 0 {

            go StartDataSaver(stor)
    }

    http.ListenAndServe(envVars.Address, r)
}
