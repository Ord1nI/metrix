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

	if serv.Config.Key != "" {
		serv.Add(middlewares.LoggerMW(serv.Logger), middlewares.SingMW(serv.Logger, []byte(serv.Config.Key)), middlewares.CompressorMW(serv.Logger))
	} else {
		serv.Add(middlewares.LoggerMW(serv.Logger), middlewares.CompressorMW(serv.Logger))
	}

	err = serv.Run()

	if err != nil {
		panic(err)
	}
}
