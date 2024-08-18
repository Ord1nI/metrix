package main

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
    _ "github.com/jackc/pgx/v5/stdlib"

    "database/sql"
	"net/http"
    "fmt"
	"time"

	"github.com/Ord1nI/metrix/internal/compressor"
	"github.com/Ord1nI/metrix/internal/logger"
	"github.com/Ord1nI/metrix/internal/storage"
    "github.com/Ord1nI/metrix/internal/configs"
)


var config configs.ServerConfig

var sugar *zap.SugaredLogger


func StartDataSaver(s *storage.MemStorage) {
    for {
        time.Sleep(time.Duration(config.StoreInterval) * time.Second)
        err := s.WriteToFile(config.FileStoragePath)
        if err != nil {
            sugar.Fatal(err)
        } else {
            sugar.Info("Data saved")
        }
    }
}

func initF() {
    log, err := logger.NewLogger()
    if err != nil {
        panic(err)
    }
    defer log.Sync()
    sugar = log.Sugar()
    sugar.Infoln("loger created successfuly")

    configs.ServerGetConf(sugar, &config)
    sugar.Info("Config vars: ", config)
}

func main() {

    initF()

    stor := storage.NewMemStorage()

    db, err := sql.Open("pgx", fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
                      `localhost`, config.Database.User, config.Database.Password, config.Database.Name))

    if err != nil {
        sugar.Fatal(err)
    }


    if config.FileStoragePath != "" && config.Restore {
        err := stor.GetFromFile(config.FileStoragePath)
        if err != nil {
            sugar.Info(err)
        } else {
            sugar.Info("Data loaded succesful",stor.Gauge)
        }
    }

    var r chi.Router

    if config.StoreInterval == 0 {
        r = CreateRouter(db, stor,
            logger.HandlerLogging(sugar), 
            compressor.GzipMiddleware(sugar), 
            storage.SaveToFileMW(sugar,config.FileStoragePath,stor))

    } else {
        r = CreateRouter(db, stor, 
            logger.HandlerLogging(sugar), 
            compressor.GzipMiddleware(sugar))
    }

    if  config.FileStoragePath != "" && 
        config.StoreInterval != 0 {
            go StartDataSaver(stor)
    }

    http.ListenAndServe(config.Address, r)
}
