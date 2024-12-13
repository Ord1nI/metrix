package main

import (
	"os"
	"os/signal"
	"syscall"

	"fmt"

	"github.com/Ord1nI/metrix/internal/server/grpcserv"
)

var (
	buildVersion string = "N/A"
	buildDate string = "N/A"
	buildCommit string = "N/A"
)

func main() {

	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)

	serv, err := grpcserv.Default()
	if err != nil {
		panic(err)

	}

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
