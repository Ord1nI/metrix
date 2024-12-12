package main

import (
	"os"
	"os/signal"
	"syscall"

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


	end := make(chan struct{})

	err = serv.Run(end)

	if err != nil {
		panic(err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-sigs

	fmt.Println("End program")

	close(end)

}
