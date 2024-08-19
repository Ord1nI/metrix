package main

import (
	"go.uber.org/zap"

	"net/http"
    "context"

	"github.com/Ord1nI/metrix/internal/compressor"
	"github.com/Ord1nI/metrix/internal/logger"
	"github.com/Ord1nI/metrix/internal/repo/storage"
	"github.com/Ord1nI/metrix/internal/repo/database"
	"github.com/Ord1nI/metrix/internal/repo"
    "github.com/Ord1nI/metrix/internal/configs"
)


var (
    config configs.ServerConfig
    
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

    configs.ServerGetConf(sugar, &config)
    sugar.Info("Config vars: ", config)
}

func initRepo() repo.Repo{
    db, err := database.NewDB(context.TODO(),config.DBdsn)
    if err != nil || db.DB.Ping() != nil {
        sugar.Error("Fail to connect to database creating memstorage")
        memStor := storage.NewMemStorage()

        if config.Restore && config.FileStoragePath != "" {
            err = memStor.GetFromFile(config.FileStoragePath)
            if err != nil {
                sugar.Infoln("unable to load data from file")
            } else {
                sugar.Infoln("successfuly load data from file")
            }
        }

        if config.StoreInterval != 0 && config.FileStoragePath != "" {
            sugar.Infoln("Starting saving to file")
            go memStor.StartDataSaver(config.StoreInterval, config.FileStoragePath) //add error check
        }
        return memStor
    } else {
        sugar.Infoln("Database loaded successfuly")
        err = db.CreateTable()
        sugar.Errorln(err)
        return db
    }
}

func main() {

    initF()

    stor := initRepo()

    defer stor.Close()

    r := CreateRouter(stor, 
        logger.HandlerLogging(sugar), 
        compressor.GzipMiddleware(sugar))

    http.ListenAndServe(config.Address, r)
}
