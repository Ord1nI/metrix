package configs

import (
	"flag"

	"github.com/Ord1nI/metrix/internal/logger"
	"github.com/caarlos0/env/v11"
)

type ServerConfig struct {
    Address string `env:"ADDRESS" envDefault:"localhost:8080"` //envvar $ADDRESS or envDefault
    StoreInterval int `env:"STORE_INTERVAL" envDefault:"300"`  //envvar $STORE_INTERVAL or envDefault
    FileStoragePath string `env:"FILE_STORAGE_PATH"`           //envvar $FILE_STORAGE or envDefault
    Restore bool `env:"RESTORE" envDefault:"true"`             //envvar $RESTORE or envDefault
} 

func ServerGetConf(sugar logger.Logger, envVars *ServerConfig) {
    err := env.Parse(envVars)

    if err != nil {
        sugar.Errorln("Couldn't get env vars")
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
    if envVars.Restore {
        envVars.Restore = *fRestore
    }
}
