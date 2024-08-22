package main

import (
    "github.com/Ord1nI/metrix/internal/server"
    "github.com/Ord1nI/metrix/internal/logger"
    "github.com/Ord1nI/metrix/internal/compressor"
)

func main() {
    serv, err := server.New()
    if err != nil {
        panic(err)
    }

    serv.Add(logger.MW(serv.Logger), compressor.MW(serv.Logger))

    serv.Init()
    defer serv.Repo.Close()

    serv.Run()
}
