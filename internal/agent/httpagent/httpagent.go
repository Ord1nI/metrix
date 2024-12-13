// Package agent contains class Agent that collect
// metrics and send them to server.
package httpagent

import (
	"github.com/go-resty/resty/v2"

	"github.com/Ord1nI/metrix/internal/agent"
)

type HttpAgent struct {
	*agent.Agent
	Client *resty.Client
}


func New() (*HttpAgent, error) {
	mAgent, err := agent.New()

	if err != nil {
		return nil, err
	}

	agent := HttpAgent{
		Agent: mAgent,
	}

	agent.Client = resty.New().SetBaseURL("http://" + agent.Config.Address)

	agent.Logger.Infoln("Agent inited successfuly")

	return &agent, nil
}

func (a *HttpAgent) Run() chan struct{} {
	end := make(chan struct{})
	if (a.Config.PublicKeyFile != "") {
		a.StartWorkers(a.TaskPoll(end, a.StartMetricCollector(end)), a.SendMetricJSONwithEncryption(a.Config.PublicKeyFile))
	} else {
		a.StartWorkers(a.TaskPoll(end, a.StartMetricCollector(end)), a.SendMetricJSON)
	}

	return end
}
