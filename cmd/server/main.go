package main

import (
	"github.com/Ord1nI/metrix/internal/compressor"
	"github.com/Ord1nI/metrix/internal/logger"
	"github.com/Ord1nI/metrix/internal/server"
	"github.com/Ord1nI/metrix/internal/signature"
)

func main() {
	serv, err := server.New()
	if err != nil {
		panic(err)
	}

    if serv.Config.Key != "" {
        serv.Add(logger.MW(serv.Logger), signature.MW(serv.Logger, []byte(serv.Config.Key)), compressor.MW(serv.Logger))
    } else {
        serv.Add(logger.MW(serv.Logger), compressor.MW(serv.Logger))
    }

	serv.Init()

	defer serv.Repo.Close()

	err = serv.Run()

	if err != nil {
		panic(err)
	}
}
