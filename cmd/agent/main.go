package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Ord1nI/metrix/internal/agent"

	"fmt"
)

var (
	buildVersion string = "N/A"
	buildDate string = "N/A"
	buildCommit string = "N/A"
)

func main() {

	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)

	agent, err := agent.New()
	if err != nil {
		panic(err)
	}

	stop := agent.Run()

	sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case <-sigs:
		fmt.Println("End program")
		close(stop);
	}

}
