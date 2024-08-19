package main

import (
	"github.com/go-resty/resty/v2"

	"github.com/Ord1nI/metrix/internal/compressor"
	"github.com/Ord1nI/metrix/internal/repo/metrics"
	"github.com/Ord1nI/metrix/internal/repo/storage"

	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"
)


func collectMetrics(stor *storage.MemStorage) {
    var mS runtime.MemStats
    runtime.ReadMemStats(&mS)
    mGauge  := storage.MGauge{
        "Alloc" : metrics.Gauge(mS.Alloc),
        "BuckHashSys" : metrics.Gauge(mS.BuckHashSys),
        "Frees" : metrics.Gauge(mS.Frees),
        "GCCPUFraction" : metrics.Gauge(mS.GCCPUFraction),
        "GCSys" : metrics.Gauge(mS.GCSys),
        "HeapAlloc" : metrics.Gauge(mS.HeapAlloc),
        "HeapIdle" : metrics.Gauge(mS.HeapIdle),
        "HeapInuse" : metrics.Gauge(mS.HeapInuse),
        "HeapObjects" : metrics.Gauge(mS.HeapObjects),
        "HeapReleased" : metrics.Gauge(mS.HeapReleased),
        "HeapSys" : metrics.Gauge(mS.HeapSys),
        "LastGC" : metrics.Gauge(mS.LastGC),
        "Lookups" : metrics.Gauge(mS.Lookups),
        "MCacheInuse" : metrics.Gauge(mS.MCacheInuse),
        "MCacheSys" : metrics.Gauge(mS.MCacheSys),
        "MSpanInuse" : metrics.Gauge(mS.MSpanInuse),
        "MSpanSys" : metrics.Gauge(mS.MSpanSys),
        "Mallocs" : metrics.Gauge(mS.Mallocs),
        "NextGC" : metrics.Gauge(mS.NextGC),
        "NumForcedGC" : metrics.Gauge(mS.NumForcedGC),
        "NumGC" : metrics.Gauge(mS.NumGC),
        "OtherSys" : metrics.Gauge(mS.OtherSys),
        "PauseTotalNs" : metrics.Gauge(mS.PauseTotalNs),
        "StackInuse" : metrics.Gauge(mS.StackInuse),
        "StackSys" : metrics.Gauge(mS.StackSys),
        "Sys" : metrics.Gauge(mS.Sys),
        "TotalAlloc" : metrics.Gauge(mS.TotalAlloc),
        "RandomValue" : metrics.Gauge(rand.Float64()),
    }

    stor.AddGauge(mGauge)

    stor.Set("PollCount",metrics.Counter(1))
}


func SendGaugeMetrics(client *resty.Client, stor *storage.MemStorage) error{
    for i, v := range *stor.Gauge {
        var builder strings.Builder
        builder.WriteString("/update/gauge/")
        builder.WriteString(i)
        builder.WriteRune('/')
        builder.WriteString(strconv.FormatFloat(float64(v), 'f', -1, 64))

        res, err := client.R().
            SetHeader("Content-Type","text/plain").
            Post(builder.String())
        
        if err != nil {
            return err
        }
        
        if res.StatusCode() != http.StatusOK {
            return errors.New("doesnt sent")
        }
    }
    return nil
}

func SendMetricsJSON(client *resty.Client, stor *storage.MemStorage) error {
    metricArr := stor.ToMetrics()

    for _, m := range metricArr {
        data, err := json.Marshal(m)
        if err != nil {
            return err
        }

        data, err = compressor.ToGzip(data)

        if err != nil {
            return err
        }


        var res *resty.Response
        for _, backoff := range envVars.BackoffSchedule {

            res, err = client.R().
                            SetHeader("Content-Type", "application/json").
                            SetHeader("Content-Encoding", "gzip").
                            SetHeader("Accept-Encoding", "gzip").
                            SetBody(data).
                            Post("/update/")

            if err == nil && res.StatusCode() == http.StatusOK{
                break
            }

            time.Sleep(backoff)
        }

        if err != nil {
            return err
        }
        
        if res.StatusCode() != http.StatusOK {
            return errors.New("doesnt sent")
        }
    }

    return nil
}

func SendMetricsArrJSON(client *resty.Client, stor *storage.MemStorage) error {
    metricsJSON, err := stor.MarshalJSON()
    if err != nil {
        sugar.Error("Error while marshaling")
        return err
    }

    metricsJSON, err = compressor.ToGzip(metricsJSON)
    if err != nil {
        sugar.Error("Error while compressing")
        return err
    }

    res, err := client.R().
                    SetHeader("Content-Type", "application/json").
                    SetHeader("Content-Encoding", "gzip").
                    SetHeader("Accept-Encoding", "gzip").
                    SetBody(metricsJSON).
                    Post("/updates/")
    if err != nil {
        sugar.Error("Error while sending request")
        return err
    }
    if res.StatusCode() != http.StatusOK {
        sugar.Infoln("get status code", res.StatusCode())
        return errors.New("StatusCode != OK")
    }
    return nil 

}

func StartClient(client *resty.Client, stor *storage.MemStorage) {
    for {
        for i := int64(0); i < envVars.ReportInterval / envVars.PollInterval; i++ {
            collectMetrics(stor)
            time.Sleep(time.Second * time.Duration(envVars.PollInterval))
        }
        err := SendMetricsArrJSON(client, stor)

        if err != nil {
            sugar.Infoln(err)
        } else {
            sugar.Infoln("Metrics sent")
        }
    }    
}
