package server

import (
	"github.com/caarlos0/env/v11"

	"flag"
	"time"
)

type Config struct {
    Address string `env:"ADDRESS" envDefault:"localhost:8080"` //envvar $ADDRESS or envDefault
    StoreInterval int `env:"STORE_INTERVAL" envDefault:"300"`  //envvar $STORE_INTERVAL or envDefault
    FileStoragePath string `env:"FILE_STORAGE_PATH"`           //envvar $FILE_STORAGE or envDefault
    Restore bool `env:"RESTORE" envDefault:"true"`             //envvar $RESTORE or envDefault
    DBdsn string `env:"DATABASE_DSN"`
    BackoffSchedule []time.Duration
} 

func (s *Server) GetConf() error{
    err := env.Parse(&s.Config)

    if err != nil {
        return err
    }

    var fAddress = flag.String("a", s.Config.Address, "enter IP format ip:port")
    var fStoreInterval = flag.Int("i", s.Config.StoreInterval,
        "enter interval (in seconds) between all data saved to specified file")
    var fFileStoragePath = flag.String("f", s.Config.FileStoragePath,
        "enter path to file where all data will be saved")
    var fRestore = flag.Bool("r", s.Config.Restore,
        "whether or not load data to specified file")
    var fDatabase = flag.String("d", s.Config.DBdsn, 
        "e.g.host=hostname user=username password=pssword dbname=dbname")

    flag.Parse()

    if s.Config.Address == "localhost:8080" {
        s.Config.Address = *fAddress
    }
    if s.Config.StoreInterval == 300 {
        s.Config.StoreInterval = *fStoreInterval
    }
    if s.Config.FileStoragePath == "" {
        s.Config.FileStoragePath = *fFileStoragePath
    }
    if s.Config.Restore {
        s.Config.Restore = *fRestore
    }
    if s.Config.DBdsn == "" {
        s.Config.DBdsn = *fDatabase
    }
    s.Config.BackoffSchedule =[]time.Duration{
        1 * time.Second,
        3 * time.Second,
        5 * time.Second,
    }
    return nil
}
