package agent

import (
	"time"

	"github.com/Ord1nI/metrix/internal/logger"
	"github.com/Ord1nI/metrix/internal/repo/storage"
	"github.com/go-resty/resty/v2"
)

type Agent struct {
    Repo *storage.MemStorage
	Logger logger.Logger
	Config Config
    Client *resty.Client
}

func New() (*Agent, error){
    log, err := logger.New()

    if err != nil {
        return nil, err
    }
    log.Infoln("Logger inited successfuly")
    Agent := Agent{
        Logger: log,
        Repo: storage.NewMemStorage(),
    }

    Agent.GetConf()
    Agent.Client = resty.New().SetBaseURL("http://" + Agent.Config.Address)

    log.Infoln("Agent inited successfuly")

    return &Agent, nil
}

func (a *Agent) Run() {
    pollTiker := time.NewTicker(time.Duration(a.Config.PollInterval) * time.Second)
    reportTicker := time.NewTicker(time.Duration(a.Config.ReportInterval) * time.Second)
    for {
        <-pollTiker.C
        a.CollectMetrics()
        a.Logger.Infoln("Metic collected")
        <-reportTicker.C
        err := a.SendMetricsArrJSON()
        if err != nil {
            a.Logger.Infoln(err)
        } else {
            a.Logger.Infoln("Metics sent")
        }
    }
}
