package main

import (
	"github.com/Ord1nI/metrix/internal/agent"
)




func main() {

    agent, err := agent.New()
    if err != nil {
        panic(err)
    }

    agent.Run()
}
