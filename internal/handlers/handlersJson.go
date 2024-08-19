package handlers

import (
	"github.com/Ord1nI/metrix/internal/repo/metrics"
	"github.com/Ord1nI/metrix/internal/repo"

	"encoding/json"
    "errors"
    "strings"
	"io"
	"net/http"
)

func UpdateJSON(s repo.GetAdder) APIFunc{
    fHandler :=  func(res http.ResponseWriter, req *http.Request) error {

        if !strings.Contains(req.Header.Get("Content-Type"), "application/json") {
            return NewHandlerError(errors.New("not json request"),http.StatusBadRequest)
        }

        data, err := io.ReadAll(req.Body)
        req.Body.Close()

        if err != nil {
            return NewHandlerError(err, http.StatusBadRequest)
        }

        var metric metrics.Metric

        err = json.Unmarshal(data, &metric)


        if err != nil {
            return NewHandlerError(err, http.StatusBadRequest)
        }

        err = s.Add(metric.ID,metric)

        if err != nil {
            return NewHandlerError(err, http.StatusBadRequest)
        }
        
        ptrMetric := metrics.Metric{
            MType: metric.MType,
        }
        s.Get(metric.ID, &ptrMetric)

        resMetric, err := json.Marshal(ptrMetric) //maybe can be error

        if err != nil {
            return NewHandlerError(err, http.StatusBadRequest)
        }

        res.Header().Add("Content-Type", "application/json" )
        res.WriteHeader(http.StatusOK)
        res.Write(resMetric)
        return nil
    }

    return APIFunc(fHandler)
}

func GetJSON(s repo.GetAdder) APIFunc {
    fHandler :=  func(res http.ResponseWriter, req *http.Request) error {

        if !strings.Contains(req.Header.Get("Content-Type"), "application/json") {
            return NewHandlerError(errors.New("not json request"),http.StatusBadRequest)
        }

        data, err := io.ReadAll(req.Body)
        req.Body.Close()

        if err != nil {
            return NewHandlerError(err, http.StatusBadRequest)
        }

        var metric metrics.Metric
        err = json.Unmarshal(data, &metric)

        if err != nil {
            return NewHandlerError(err, http.StatusBadRequest)
        }

        err = s.Get(metric.ID, &metric)

        if err != nil {
            return NewHandlerError(err, http.StatusNotFound)
        }

        resMetric, err := json.Marshal(metric) //maybe can be error

        if err != nil {
            return NewHandlerError(err, http.StatusBadRequest)
        }

        res.Header().Add("Content-Type", "application/json" )
        res.WriteHeader(http.StatusOK)
        res.Write(resMetric)
        return nil
    }
    return APIFunc(fHandler)
}

func UpdatesJSON(s repo.Repo) APIFunc { 
    fHandler :=  func(res http.ResponseWriter, req *http.Request) error {

        if !strings.Contains(req.Header.Get("Content-Type"), "application/json") {
            return NewHandlerError(errors.New("not json request"),http.StatusBadRequest)
        }
        
        data, err := io.ReadAll(req.Body)
        req.Body.Close()

        if err != nil {
            return NewHandlerError(err, http.StatusBadRequest)
        }

        var metrics []metrics.Metric
        err = json.Unmarshal(data, &metrics)
        
        if err != nil {
            return NewHandlerError(err, http.StatusBadRequest)
        }

        err = s.Add("", metrics)
        if err != nil {
            return NewHandlerError(err, http.StatusBadRequest)
        }

        res.Header().Add("Content-Type", "application/json" )
        res.WriteHeader(http.StatusOK)
        res.Write([]byte("metrics added"))
        return nil
    }
    return APIFunc(fHandler)
}
