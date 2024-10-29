package agent

import (
	"time"

	"github.com/Ord1nI/metrix/internal/repo/metrics"
)

func (a *Agent) StartMetricCollector(stop <-chan struct{}) chan struct{} {
	OK := make(chan struct{})
	pollTiker := time.NewTicker(time.Duration(a.Config.PollInterval) * time.Second)
	go func() {
		for {
			select {
			case <-stop:
				a.Logger.Infoln("stop metricCollercor gorutine")
				close(OK)
				return
			default:
				<-pollTiker.C
				a.CollectMetrics()
				OK <- struct{}{}
			}
		}
	}()
	return OK
}

func (a *Agent) TaskPoll(stop <-chan struct{}, ok chan struct{}) chan metrics.Metric {
	taskPoll := make(chan metrics.Metric)
	ReportTiker := time.NewTicker(time.Duration(a.Config.ReportInterval) * time.Second)

	go func() {
		for {
			select {
			case <-stop:
				a.Logger.Infoln("stop taskPool gorutine")
				close(taskPoll)
				return
			default:
				<-ReportTiker.C
				<-ok
				var taskList []metrics.Metric

				if a.Repo.Get("", &taskList) != nil {
					a.Logger.Fatal("Failed to init task list")
				}

				for _, m := range taskList {
					taskPoll <- m
				}
			}
		}
	}()
	return taskPoll
}

func (a *Agent) StartWorkers(jobs <-chan metrics.Metric, sendFunc func(metrics.Metric)error) {
	for i := range a.Config.RateLimit {
		a.Logger.Infoln("start", i, "worker")
		go func() {
			for j := range jobs {
				err := sendFunc(j)
				if err != nil {
					a.Logger.Infoln(err)
				}
			}
			a.Logger.Infoln("End worker")
		}()
	}
}
