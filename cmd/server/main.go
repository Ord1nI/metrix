package main

import (
    "github.com/caarlos0/env/v11"

    "flag"
    "net/http"
    "github.com/Ord1nI/metrix/internal/storage"
    "github.com/Ord1nI/metrix/internal/handlers"
)
type Config struct {
    Address string `env:"ADDRESS"`
}

var envVars Config
var fIPStr *string 

func init() {

    env.Parse(&envVars)

    if envVars.Address == "" {
        fIPStr = flag.String("a", "localhost:8080", "enter IP format ip:port")
    } else {
        fIPStr = &envVars.Address
    }
}

func main() {
    flag.Parse()

    stor := storage.NewEmptyStorage()

    r := CreateRouter(stor)

    r.Get("/", handlers.GetAllMetrics(stor))                  //POST localhost:/


    err := http.ListenAndServe(*fIPStr, r)
    if err != nil {
        panic(err)
    }
}
