package main

import (
	"github.com/Ord1nI/metrix/internal/middlewares"
	"github.com/Ord1nI/metrix/internal/server"
)

func main() {
    serv, err := server.New()
    if err != nil {
        panic(err)
    }

    serv.Add(middlewares.LoggerMW(serv.Logger), middlewares.CompressorMW(serv.Logger))

    err = serv.Run()

    if err != nil{
        panic(err)
    }
}
