//Package agent contains class Agent that collect
//metrics and send them to server.
package agent

import (
	"github.com/go-resty/resty/v2"

	"github.com/Ord1nI/metrix/internal/logger"
	"github.com/Ord1nI/metrix/internal/repo/storage"
)

type Agent struct {
	Logger logger.Logger
	Repo   *storage.MemStorage
	Client *resty.Client
	Config Config
}

func New() (*Agent, error) {
	log, err := logger.New()

	if err != nil {
		return nil, err
	}

	log.Infoln("Logger inited successfuly")
	agent := Agent{
		Logger: log,
		Repo:   storage.NewMemStorage(),
	}

	agent.GetConf()
	agent.Client = resty.New().SetBaseURL("http://" + agent.Config.Address)

	log.Infoln("Agent inited successfuly")

	return &agent, nil
}

func (a *Agent) Run() chan struct{} {
	end := make(chan struct{})
	a.StartWorkers(a.TaskPoll(end, a.StartMetricCollector(end)))
	return end
}
