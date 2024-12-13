package server

import (
	"encoding/json"
	"io"
	"os"

	"github.com/caarlos0/env/v11"

	"flag"
	"time"
)

type Config struct {
	Address         string `json:"address" env:"ADDRESS" envDefault:"localhost:8080"`
	FileStoragePath string `json:"file_storrage_path" env:"FILE_STORAGE_PATH"`
	DBdsn           string `json:"database_dsn" env:"DATABASE_DSN"`
	Key             string `json:"key" env:"KEY" envDefault:""`
	FileCfg         string `env:"SERVER_CONFIG" envDefault:""`
	PrivateKeyFile  string `json:"crypto_key" env:"CRYPTO_KEY" envDefault:""`
	TrustedSubnet   string `json:"trusted_subnet" env:"TRUSTED_SUBNET" envDefault:""`
	BackoffSchedule []time.Duration
	StoreInterval   int64  `json:"store_interval" env:"STORE_INTERVAL" envDefault:"300"`
	Restore         bool `json:"restore" env:"RESTORE" envDefault:"true"`
}

func setValue[i int64 | string | bool] (mainCfg *i, fileCfg *i, flag *i, defaultValue i) {
	if *mainCfg == defaultValue {
		*mainCfg = *flag;
	}

	var sValue i;

	if *mainCfg == defaultValue && *fileCfg != sValue{
		*mainCfg = *fileCfg
	}
}

func (s *Server) getConfFromFile(fileName string) Config {



	file, err := os.OpenFile(fileName, os.O_RDONLY, 0644);

	if err != nil {
		s.Logger.Fatalln(err);
	}

	var cfg = Config{}

	bFile, err := io.ReadAll(file)

	if err != nil {
		s.Logger.Fatalln(err)
	}

	err = json.Unmarshal(bFile, &cfg)

	if err != nil {
		s.Logger.Fatalln(err)
	}

	return cfg
}

func (s *Server) GetConf() error {
	err := env.Parse(s.Config)

	if err != nil {
		return err
	}

	var (
		fAddress       = flag.String("a", s.Config.Address, "enter IP format ip:port")
		fStoreInterval = flag.Int64("i", s.Config.StoreInterval,
			"enter interval (in seconds) between all data saved to specified file")
		fFileStoragePath = flag.String("f", s.Config.FileStoragePath,
			"enter path to file where all data will be saved")
		fRestore = flag.Bool("r", s.Config.Restore,
			"whether or not load data to specified file")
		fFileCfg = flag.String("config", s.Config.FileCfg, "enter config file location")
		fDatabase = flag.String("d", s.Config.DBdsn,
			"e.g.host=hostname user=username password=pssword dbname=dbname")
		fKey = flag.String("k", s.Config.Key, "enter Signatur key")

		fPrivateKeyFile = flag.String("crypto-key", s.Config.PrivateKeyFile, "enter location of file with private key")

		fTrustedSubnet = flag.String("t", s.Config.TrustedSubnet, "enter from what subnet you want to recieve requests")
	)

	flag.Parse()

	if s.Config.FileCfg == "" {
		s.Config.FileCfg = *fFileCfg
	}

	var cfgFromFile Config

	if s.Config.FileCfg != ""  {
		cfgFromFile = s.getConfFromFile(s.Config.FileCfg)
	}

	setValue(&s.Config.Address, &cfgFromFile.Address, fAddress, "localhost:8080")
	setValue(&s.Config.StoreInterval, &cfgFromFile.StoreInterval, fStoreInterval, 300)
	setValue(&s.Config.FileStoragePath, &cfgFromFile.FileStoragePath, fFileStoragePath, "")
	setValue(&s.Config.Restore, &cfgFromFile.Restore, fRestore, true)
	setValue(&s.Config.DBdsn, &cfgFromFile.DBdsn, fDatabase, "")
	setValue(&s.Config.Key, &cfgFromFile.Key, fKey, "")
	setValue(&s.Config.PrivateKeyFile, &cfgFromFile.PrivateKeyFile, fPrivateKeyFile, "")
	setValue(&s.Config.TrustedSubnet, &cfgFromFile.TrustedSubnet, fTrustedSubnet, "")

	s.Config.BackoffSchedule = []time.Duration{
		1 * time.Second,
		3 * time.Second,
		5 * time.Second,
	}
	return nil
}
