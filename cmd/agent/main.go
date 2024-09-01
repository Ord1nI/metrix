package main

import (
	"time"

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
    time.Sleep(time.Second * 150)
}
