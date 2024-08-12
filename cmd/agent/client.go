package main

import (
	"github.com/go-resty/resty/v2"

	"github.com/Ord1nI/metrix/internal/myjson"
	"github.com/Ord1nI/metrix/internal/storage"

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
        "Alloc" : storage.Gauge(mS.Alloc),
        "BuckHashSys" : storage.Gauge(mS.BuckHashSys),
        "Frees" : storage.Gauge(mS.Frees),
        "GCCPUFraction" : storage.Gauge(mS.GCCPUFraction),
        "GCSys" : storage.Gauge(mS.GCSys),
        "HeapAlloc" : storage.Gauge(mS.HeapAlloc),
        "HeapIdle" : storage.Gauge(mS.HeapIdle),
        "HeapInuse" : storage.Gauge(mS.HeapInuse),
        "HeapObjects" : storage.Gauge(mS.HeapObjects),
        "HeapReleased" : storage.Gauge(mS.HeapReleased),
        "HeapSys" : storage.Gauge(mS.HeapSys),
        "LastGC" : storage.Gauge(mS.LastGC),
        "Lookups" : storage.Gauge(mS.Lookups),
        "MCacheInuse" : storage.Gauge(mS.MCacheInuse),
        "MCacheSys" : storage.Gauge(mS.MCacheSys),
        "MSpanInuse" : storage.Gauge(mS.MSpanInuse),
        "MSpanSys" : storage.Gauge(mS.MSpanSys),
        "Mallocs" : storage.Gauge(mS.Mallocs),
        "NextGC" : storage.Gauge(mS.NextGC),
        "NumForcedGC" : storage.Gauge(mS.NumForcedGC),
        "NumGC" : storage.Gauge(mS.NumGC),
        "OtherSys" : storage.Gauge(mS.OtherSys),
        "PauseTotalNs" : storage.Gauge(mS.PauseTotalNs),
        "StackInuse" : storage.Gauge(mS.StackInuse),
        "StackSys" : storage.Gauge(mS.StackSys),
        "Sys" : storage.Gauge(mS.Sys),
        "TotalAlloc" : storage.Gauge(mS.TotalAlloc),
        "RandomValue" : storage.Gauge(rand.Float64()),
    }

    stor.AddGauge(mGauge)

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
    metrics := stor.ToMetrics()
    delta := int64(1)
    metrics = append(metrics, myjson.Metric{ID:"PollCount",MType: "counter", Delta: &delta})

    for _, m := range metrics {
        data, err := json.Marshal(m)

        if err != nil {
            return errors.New("marshal error")
        }
        backoffSchedule := []time.Duration{
          100 * time.Millisecond,
          500 * time.Millisecond,
          1 * time.Second,
        }

        var res *resty.Response
        for _, backoff := range backoffSchedule {

            res, err = client.R().
                            SetHeader("Content-Type", "application/json").
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


func StartClient(client *resty.Client, stor *storage.MemStorage) {
    for {
        for i := int64(0); i < envVars.ReportInterval / envVars.PollInterval; i++ {
            collectMetrics(stor)
            time.Sleep(time.Second * time.Duration(envVars.PollInterval))
        }
        err := SendMetricsJSON(client, stor)

        if err != nil {
            sugar.Infoln(err)
        } else {
            sugar.Infoln("Metrics sent")
        }
    }    
}
