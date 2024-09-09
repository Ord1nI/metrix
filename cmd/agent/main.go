package main

import (
	"github.com/Ord1nI/metrix/internal/agent"
)

func main() {

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
