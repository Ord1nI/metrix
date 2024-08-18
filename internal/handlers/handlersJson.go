package handlers

import (
	"github.com/Ord1nI/metrix/internal/repo/metrics"
	"github.com/Ord1nI/metrix/internal/logger"
	"github.com/Ord1nI/metrix/internal/repo"

	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func UpdateJSON(l logger.Logger, s repo.MetricGetAdder) http.Handler{
    fHandler :=  func(res http.ResponseWriter, req *http.Request) {

        if req.Header.Get("Content-Type") != "application/json" {
            res.WriteHeader(http.StatusBadRequest)
            res.Write([]byte("not json request\n"))
            return
        }

        data, err := io.ReadAll(req.Body)
        req.Body.Close()

        if err != nil {
            l.Errorln(err)
            res.WriteHeader(http.StatusBadRequest)
            res.Write([]byte("Bad rquest body\n"))
            return
        }

        var metric metrics.Metric

        err = json.Unmarshal(data, &metric)


        if err != nil {
            l.Errorln(err)
            res.WriteHeader(http.StatusBadRequest)
            res.Write([]byte("Cant unmarshal json\n"))
            return
        }

        err = s.AddMetric(metric)

        if err != nil {
            l.Errorln(err)
            res.WriteHeader(http.StatusBadRequest)
            res.Write([]byte(fmt.Sprint(err)))
            return
        }
        
        ptrMetric, _ := s.GetMetric(metric.ID, metric.MType)

        resMetric, _ := json.Marshal(ptrMetric) //maybe can be error

        res.Header().Add("Content-Type", "application/json" )
        res.WriteHeader(http.StatusOK)
        res.Write(resMetric)
    }

    return http.HandlerFunc(fHandler)
}

func GetJSON(l logger.Logger, s repo.MetricGetAdder) http.Handler {
    fHandler :=  func(res http.ResponseWriter, req *http.Request) {

        if req.Header.Get("Content-Type") != "application/json" {
            res.WriteHeader(http.StatusBadRequest)
            res.Write([]byte("not json request\n"))
            return
        }

        data, err := io.ReadAll(req.Body)
        req.Body.Close()

        if err != nil {
            l.Errorln(err)
            res.WriteHeader(http.StatusBadRequest)
            res.Write([]byte("Bad rquest body\n"))
            return
        }

        var metric metrics.Metric
        err = json.Unmarshal(data, &metric)

        if err != nil {
            l.Errorln(err)
            res.WriteHeader(http.StatusBadRequest)
            res.Write([]byte("Cant unmarshal json\n"))
            return
        }

        ptrMetric, ok := s.GetMetric(metric.ID, metric.MType)

        if !ok {
            res.WriteHeader(http.StatusNotFound)
            res.Write([]byte("Cant find this metric"))
            return
        }

        resMetric, _ := json.Marshal(ptrMetric) //maybe can be error

        res.Header().Add("Content-Type", "application/json" )
        res.WriteHeader(http.StatusOK)
        res.Write(resMetric)
    }
    return http.HandlerFunc(fHandler)
}
