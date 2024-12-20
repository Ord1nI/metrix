package server

import (
	"github.com/caarlos0/env/v11"

	"flag"
	"time"
)

type Config struct {
	Address         string `env:"ADDRESS" envDefault:"localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DBdsn           string `env:"DATABASE_DSN"`
	Key             string `env:"KEY" envDefault:""`
	BackoffSchedule []time.Duration
	StoreInterval   int  `env:"STORE_INTERVAL" envDefault:"300"`
	Restore         bool `env:"RESTORE" envDefault:"true"`
}

func (s *Server) GetConf() error {
	err := env.Parse(&s.Config)

	if err != nil {
		return err
	}

	var (
		fAddress       = flag.String("a", s.Config.Address, "enter IP format ip:port")
		fStoreInterval = flag.Int("i", s.Config.StoreInterval,
			"enter interval (in seconds) between all data saved to specified file")
		fFileStoragePath = flag.String("f", s.Config.FileStoragePath,
			"enter path to file where all data will be saved")
		fRestore = flag.Bool("r", s.Config.Restore,
			"whether or not load data to specified file")
		fDatabase = flag.String("d", s.Config.DBdsn,
			"e.g.host=hostname user=username password=pssword dbname=dbname")
		fKey = flag.String("k", s.Config.Key, "enter Signatur key")
	)

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
	if s.Config.Key == "" {
		s.Config.Key = *fKey
	}

	s.Config.BackoffSchedule = []time.Duration{
		1 * time.Second,
		3 * time.Second,
		5 * time.Second,
	}
	return nil
}
