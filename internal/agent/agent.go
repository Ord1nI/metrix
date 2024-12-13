package agent

import (
	"fmt"
	"runtime"

	mrand "math/rand/v2"

	"github.com/Ord1nI/metrix/internal/logger"
	"github.com/Ord1nI/metrix/internal/repo/metrics"
	"github.com/Ord1nI/metrix/internal/repo/storage"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

type Agent struct{
	Logger logger.Logger
	Repo   *storage.MemStorage
	Config *Config
}

func (a *Agent) CollectMetrics() {
	var mS runtime.MemStats
	runtime.ReadMemStats(&mS)
	memory, _ := mem.VirtualMemory()
	mGauge := storage.MGauge{

		"Alloc":         metrics.Gauge(mS.Alloc),
		"BuckHashSys":   metrics.Gauge(mS.BuckHashSys),
		"Frees":         metrics.Gauge(mS.Frees),
		"GCCPUFraction": metrics.Gauge(mS.GCCPUFraction),
		"GCSys":         metrics.Gauge(mS.GCSys),
		"HeapAlloc":     metrics.Gauge(mS.HeapAlloc),
		"HeapIdle":      metrics.Gauge(mS.HeapIdle),
		"HeapInuse":     metrics.Gauge(mS.HeapInuse),
		"HeapObjects":   metrics.Gauge(mS.HeapObjects),
		"HeapReleased":  metrics.Gauge(mS.HeapReleased),
		"HeapSys":       metrics.Gauge(mS.HeapSys),
		"LastGC":        metrics.Gauge(mS.LastGC),
		"Lookups":       metrics.Gauge(mS.Lookups),
		"MCacheInuse":   metrics.Gauge(mS.MCacheInuse),
		"MCacheSys":     metrics.Gauge(mS.MCacheSys),
		"MSpanInuse":    metrics.Gauge(mS.MSpanInuse),
		"MSpanSys":      metrics.Gauge(mS.MSpanSys),
		"Mallocs":       metrics.Gauge(mS.Mallocs),
		"NextGC":        metrics.Gauge(mS.NextGC),
		"NumForcedGC":   metrics.Gauge(mS.NumForcedGC),
		"NumGC":         metrics.Gauge(mS.NumGC),
		"OtherSys":      metrics.Gauge(mS.OtherSys),
		"PauseTotalNs":  metrics.Gauge(mS.PauseTotalNs),
		"StackInuse":    metrics.Gauge(mS.StackInuse),
		"StackSys":      metrics.Gauge(mS.StackSys),
		"Sys":           metrics.Gauge(mS.Sys),
		"TotalAlloc":    metrics.Gauge(mS.TotalAlloc),
		"RandomValue":   metrics.Gauge(mrand.Float64()),

		"TotalMemory": metrics.Gauge(memory.Total),
		"FreeMemory":  metrics.Gauge(memory.Free),
	}
	cpuUtil, _ := cpu.Percent(0, true)

	for i, v := range cpuUtil {
		mGauge.Add(fmt.Sprintf("CPUutilization%d", i+1), metrics.Gauge(v))
	}

	a.Repo.AddGauge(mGauge)

	a.Repo.Set("PollCount", metrics.Counter(1))
}

func New() (*Agent, error) {
	log, err := logger.New()

	if err != nil {
		return nil, err
	}

	log.Infoln("Logger inited successfuly")
	agent := &Agent{
		Logger: log,
		Repo:   storage.NewMemStorage(),
		Config: &Config{},
	}

	agent.GetConf()

	return agent, nil
}
