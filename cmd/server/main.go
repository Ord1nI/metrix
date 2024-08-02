package main

import (
    "github.com/caarlos0/env/v11"

    "flag"
    "net/http"
    "github.com/Ord1nI/metrix/internal/storage"
)

type Config struct {
    Address string `env:"ADDRESS" envDefault:"localhost:8080"` //envvar $ADDRESS or envDefault
}

var envVars Config

func getConf() {
    err := env.Parse(&envVars)

    if err != nil {
        panic(err)
    }

    fAddress := flag.String("a", envVars.Address, "enter IP format ip:port")

    flag.Parse()

    if envVars.Address == "localhost:8080" {
        envVars.Address = *fAddress
    }
}

func init() {
    getConf()
}

func main() {

    stor := storage.NewEmptyStorage()

    r := CreateRouter(stor)


    err := http.ListenAndServe(envVars.Address, r)
    if err != nil {
        panic(err)
    }
}
