package main

import (
	"github.com/go-chi/chi/v5"
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


var config configs.ServerConfig

var sugar *zap.SugaredLogger


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
    if err != nil {
        sugar.Error("Fail to connect to database creating memstorage")
        memStor := storage.NewMemStorage()

        if config.StoreInterval != 0 && config.FileStoragePath != "" {
            go memStor.StartDataSaver(config.StoreInterval, config.FileStoragePath)
            if err != nil {
                sugar.Error("fail to enable file saving")
            } else {
                sugar.Info("Successfuly enable file saving")
            }
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



    var r chi.Router

    r = CreateRouter(stor, 
        logger.HandlerLogging(sugar), 
        compressor.GzipMiddleware(sugar))

    http.ListenAndServe(config.Address, r)
}
