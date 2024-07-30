package main

import (
    "github.com/go-resty/resty/v2"

    "runtime"
    "strconv"
    "math/rand"
    "strings"
    "errors"
    "time"
    "fmt"
    "net/http"
)


func collectMetrics() {
    var mS runtime.MemStats
    runtime.ReadMemStats(&mS)
    metrics = map[string]float64{
        "Alloc" : float64(mS.Alloc),
        "BuckHashSys" : float64(mS.BuckHashSys),
        "Frees" : float64(mS.Frees),
        "GCCPUFraction" : float64(mS.GCCPUFraction),
        "GCSys" : float64(mS.GCSys),
        "HeapAlloc" : float64(mS.HeapAlloc),
        "HeapIdle" : float64(mS.HeapIdle),
        "HeapInuse" : float64(mS.HeapInuse),
        "HeapObjects" : float64(mS.HeapObjects),
        "HeapReleased" : float64(mS.HeapReleased),
        "HeapSys" : float64(mS.HeapSys),
        "LastGC" : float64(mS.LastGC),
        "Lookups" : float64(mS.Lookups),
        "MCacheInuse" : float64(mS.MCacheInuse),
        "MCacheSys" : float64(mS.MCacheSys),
        "MSpanInuse" : float64(mS.MSpanInuse),
        "MSpanSys" : float64(mS.MSpanSys),
        "Mallocs" : float64(mS.Mallocs),
        "NextGC" : float64(mS.NextGC),
        "NumForcedGC" : float64(mS.NumForcedGC),
        "NumGC" : float64(mS.NumGC),
        "OtherSys" : float64(mS.OtherSys),
        "PauseTotalNs" : float64(mS.PauseTotalNs),
        "StackInuse" : float64(mS.StackInuse),
        "StackSys" : float64(mS.StackSys),
        "Sys" : float64(mS.Sys),
        "TotalAlloc" : float64(mS.TotalAlloc),
        "RandomValue" : rand.Float64(),
    }
}


func SendGaugeMetrics(client *resty.Client) error{
    for i, v := range metrics {
        var builder strings.Builder
        builder.WriteString("/update/gauge/")
        builder.WriteString(i)
        builder.WriteRune('/')
        builder.WriteString(strconv.FormatFloat(v, 'f', -1, 64))
        res, err := client.R().
            ExpectContentType("text/plain").
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
func StartClient(client *resty.Client) {
    reportInterval := *fReportInterval
    pollInterval := *fPollInterval
    for {
        for i := int64(0); i < reportInterval / pollInterval; i++ {
            collectMetrics()
            time.Sleep(time.Second * time.Duration(pollInterval))
        }

        err := SendGaugeMetrics(client)

        if err != nil {
            fmt.Println(err)
        } else {
            fmt.Println("Gauge metrics sent")
        }

        res, err := client.R().
            ExpectContentType("text/plain").
            Post("/update/counter/PollCount/1")

        if err != nil || res.StatusCode() != http.StatusOK{
            fmt.Println("Counter metrics wasnt't sended")
        } else {
            fmt.Println("Counter metrics sented")
        }
    }    
}
