package handlers

import (
	"github.com/Ord1nI/metrix/internal/repo/metrics"
	"github.com/Ord1nI/metrix/internal/logger"
	"github.com/Ord1nI/metrix/internal/repo"

	"encoding/json"
    "strings"
	"io"
	"net/http"
)

func UpdateJSON(l logger.Logger, s repo.MetricGetAdder) http.Handler{
    fHandler :=  func(res http.ResponseWriter, req *http.Request) {

        if !strings.Contains(req.Header.Get("Content-Type"), "application/json") {
            l.Infoln("Content Type doesn't contains application/json")
            http.Error(res, "Not json request", http.StatusBadRequest)
            return
        }

        data, err := io.ReadAll(req.Body)
        req.Body.Close()

        if err != nil {
            l.Infoln(err)
            http.Error(res, "Error while updating", http.StatusBadRequest)
            return
        }

        var metric metrics.Metric

        err = json.Unmarshal(data, &metric)


        if err != nil {
            l.Infoln(err)
            http.Error(res, "Error while updating", http.StatusBadRequest)
            return
        }

        err = s.AddMetric(metric)

        if err != nil {
            l.Infoln(err)
            http.Error(res, "Error while updating", http.StatusBadRequest)
            return
        }
        
        ptrMetric, _ := s.GetMetric(metric.ID, metric.MType)

        resMetric, err := json.Marshal(ptrMetric) //maybe can be error

        if err != nil {
            l.Infoln(err)
            http.Error(res, "Error while updating", http.StatusBadRequest)
            return
        }

        res.Header().Add("Content-Type", "application/json" )
        res.WriteHeader(http.StatusOK)
        res.Write(resMetric)
    }

    return http.HandlerFunc(fHandler)
}

func GetJSON(l logger.Logger, s repo.MetricGetAdder) http.Handler {
    fHandler :=  func(res http.ResponseWriter, req *http.Request) {

        if !strings.Contains(req.Header.Get("Content-Type"), "application/json") {
            l.Infoln("Content Type doesn't contains application/json")
            http.Error(res, "Not json request", http.StatusBadRequest)
            return
        }

        data, err := io.ReadAll(req.Body)
        req.Body.Close()

        if err != nil {
            l.Infoln(err)
            http.Error(res, "Error while getting", http.StatusBadRequest)
            return
        }

        var metric metrics.Metric
        err = json.Unmarshal(data, &metric)

        if err != nil {
            l.Infoln(err)
            http.Error(res, "Error while gettings", http.StatusBadRequest)
            return
        }

        ptrMetric, ok := s.GetMetric(metric.ID, metric.MType)

        if !ok {
            l.Infoln("Metric not found")
            http.Error(res, "Metric not found", http.StatusNotFound)
            return
        }

        resMetric, err := json.Marshal(ptrMetric) //maybe can be error

        if err != nil {
            l.Infoln(err)
            http.Error(res, "Error while gettings", http.StatusBadRequest)
            return
        }

        res.Header().Add("Content-Type", "application/json" )
        res.WriteHeader(http.StatusOK)
        res.Write(resMetric)
    }
    return http.HandlerFunc(fHandler)
}

func UpdatesJSON(l logger.Logger, s repo.MetricGetAdder) http.Handler{
    fHandler :=  func(res http.ResponseWriter, req *http.Request) {
        if !strings.Contains(req.Header.Get("Content-Type"), "application/json") {
            l.Infoln("Content Type doesn't contains application/json")
            http.Error(res, "Not json request", http.StatusBadRequest)
            return
        }
        
        data, err := io.ReadAll(req.Body)
        req.Body.Close()

        if err != nil {
            l.Infoln(err)
            http.Error(res, "Error while getting", http.StatusBadRequest)
            return
        }

        var metrics []metrics.Metric
        err = json.Unmarshal(data, &metrics)
        
        if err != nil {
            l.Infoln(err)
            http.Error(res, "Error while updating", http.StatusNotFound)
            return
        }

        for _, m := range metrics {
            err = s.AddMetric(m)
            if err != nil {
                l.Infoln(err)
                http.Error(res, "Error while updating", http.StatusNotFound)
                return
            }
        }

        res.Header().Add("Content-Type", "application/json" )
        res.WriteHeader(http.StatusOK)
        res.Write([]byte("metrics added"))
    }
    return http.HandlerFunc(fHandler)
}
