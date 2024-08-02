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

func init() {

    env.Parse(&envVars)

    if envVars.Address == "" {
        envVars.Address = *flag.String("a", "localhost:8080", "enter IP format ip:port")
    }
}

func main() {
    flag.Parse()

    stor := storage.NewEmptyStorage()

    r := CreateRouter(stor)


    err := http.ListenAndServe(envVars.Address, r)
    if err != nil {
        panic(err)
    }
}
