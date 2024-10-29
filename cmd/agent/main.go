package main

import (
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
	if stop != nil {
		defer close(stop)
	}
	<-stop
}
