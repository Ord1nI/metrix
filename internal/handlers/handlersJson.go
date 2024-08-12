package handlers

import (
    "github.com/Ord1nI/metrix/internal/storage"
    "github.com/Ord1nI/metrix/internal/myjson"

    "io"
    "fmt"
    "net/http"
    "encoding/json"
)

func UpdateJSON(s storage.MetricGetAdder) http.Handler{
    fHandler :=  func(res http.ResponseWriter, req *http.Request) {

        if req.Header.Get("Content-Type") != "application/json" {
            res.WriteHeader(http.StatusBadRequest)
            res.Write([]byte("not json request\n"))
        }

        data, err := io.ReadAll(req.Body)
        req.Body.Close()

        if err != nil {
            res.WriteHeader(http.StatusBadRequest)
            res.Write([]byte("Bad rquest body\n"))
        }

        var metric myjson.Metric
        err = json.Unmarshal(data, &metric)

        fmt.Println(metric)

        if err != nil {
            res.WriteHeader(http.StatusBadRequest)
            res.Write([]byte("Cant unmarshal json\n"))
        }

        s.AddMetric(metric)
        
        ptrMetric, _ := s.GetMetric(metric.ID, metric.MType)

        resMetric, _ := json.Marshal(ptrMetric) //maybe can be error

        res.WriteHeader(http.StatusOK)
        res.Write(resMetric)

    }

    return http.HandlerFunc(fHandler)
}

func GetJSON(s storage.MetricGetAdder) http.Handler {
    fHandler :=  func(res http.ResponseWriter, req *http.Request) {
        if req.Header.Get("Content-Type") != "application/json" {
            res.WriteHeader(http.StatusBadRequest)
            res.Write([]byte("not json request\n"))
        }

        data, err := io.ReadAll(req.Body)
        req.Body.Close()

        if err != nil {
            res.WriteHeader(http.StatusBadRequest)
            res.Write([]byte("Bad rquest body\n"))
        }

        var metric myjson.Metric
        err = json.Unmarshal(data, &metric)

        if err != nil {
            res.WriteHeader(http.StatusBadRequest)
            res.Write([]byte("Cant unmarshal json\n"))
        }

        ptrMetric, ok := s.GetMetric(metric.ID, metric.MType)

        if !ok {
            res.WriteHeader(http.StatusNotFound)
            res.Write([]byte("No metric with this name"))
        }

        resMetric, _ := json.Marshal(ptrMetric) //maybe can be error

        res.WriteHeader(http.StatusOK)
        res.Write(resMetric)
    }
    return http.HandlerFunc(fHandler)
}


