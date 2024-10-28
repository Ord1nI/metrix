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

	serv.Add(middlewares.LoggerMW(serv.Logger))

	if serv.Config.PrivateKeyFile != "" {
		serv.Add(middlewares.Decrypt(serv.Logger,serv.Config.PrivateKeyFile))
	}
	if serv.Config.Key != "" {
		serv.Add(middlewares.SingMW(serv.Logger, []byte(serv.Config.Key)))
	}

	serv.Add(middlewares.CompressorMW(serv.Logger))


	err = serv.Run()

	if err != nil {
		panic(err)
	}
}
