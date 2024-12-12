package main

import (
	"github.com/Ord1nI/metrix/internal/middlewares"
	"github.com/Ord1nI/metrix/internal/server"

	"fmt"
)

var (
	buildVersion string = "N/A"
	buildDate string = "N/A"
	buildCommit string = "N/A"
)

func main() {

	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)

	serv, err := server.New()
	if err != nil {
		panic(err)
	}

	serv.Add(middlewares.LoggerMW(serv.Logger))

	if serv.Config.PrivateKeyFile != "" {
		serv.Add(middlewares.Decrypt(serv.Logger,serv.Config.PrivateKeyFile))
	}

	if serv.Config.Key != "" {
		serv.Add(middlewares.SignMW(serv.Logger, []byte(serv.Config.Key)))
	}

	serv.Add(middlewares.CompressorMW(serv.Logger))


	err = serv.Run()

	if err != nil {
		panic(err)
	}
}
