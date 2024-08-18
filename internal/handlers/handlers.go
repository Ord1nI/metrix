package handlers

import (
	"github.com/go-chi/chi/v5"

	"bytes"
	"cmp"
	"encoding/json"
	"io"
	"net/http"
	"slices"
	"strconv"

	"github.com/Ord1nI/metrix/internal/repo/metrics"
	"github.com/Ord1nI/metrix/internal/repo"
	"github.com/Ord1nI/metrix/internal/logger"
)


func UpdateGauge(l logger.Logger,s repo.Adder) http.Handler {
    fHandler := func(res http.ResponseWriter, req *http.Request) {
        name := chi.URLParam(req, "name")
        v := chi.URLParam(req, "val")

        val,err := strconv.ParseFloat(v, 64)

        if err != nil {
            l.Infoln(err)
            http.Error(res, "Error while updating", http.StatusBadRequest)
            return
        }

        s.Add(name, metrics.Gauge(val))
        res.WriteHeader(http.StatusOK)
    }
    return http.HandlerFunc(fHandler)
}

func UpdateCounter(l logger.Logger, s repo.Adder) http.Handler{
    fHandler := func(res http.ResponseWriter, req *http.Request) {

        name := chi.URLParam(req, "name")
        v := chi.URLParam(req, "val")
        
        val, err := strconv.ParseInt(v, 10, 64)

        if err != nil {
            l.Infoln(err)
            http.Error(res, "Error while updating", http.StatusBadRequest)
            return
        }

        s.Add(name, metrics.Counter(val))
        res.WriteHeader(http.StatusOK)
    }
    return http.HandlerFunc(fHandler)
}

func GetGauge(l logger.Logger,s repo.Getter) http.Handler {
    fHandler :=  func(res http.ResponseWriter, req *http.Request) {
        name := chi.URLParam(req,"name")
        var v metrics.Gauge
        err := s.Get(name, &v)

        if err != nil {
            l.Infoln(err)
            http.Error(res, "Metric not found", http.StatusNotFound)
            return
        }

        res.WriteHeader(http.StatusOK)
        io.WriteString(res, strconv.FormatFloat(float64(v), 'f', -1, 64))
        res.Write([]byte("\n"))
    }
    return http.HandlerFunc(fHandler)
}

func GetCounter(l logger.Logger, s repo.Getter) http.Handler {

    fHandler :=  func(res http.ResponseWriter, req *http.Request) {
        name := chi.URLParam(req,"name")
        var v metrics.Counter
        err := s.Get(name, &v)

        if err != nil {
            l.Infoln(err)
            http.Error(res, "Metric not found", http.StatusNotFound)
            return
        }

        res.WriteHeader(http.StatusOK)
        io.WriteString(res, strconv.FormatInt(int64(v), 10))
        res.Write([]byte("\n"))
    }
    return http.HandlerFunc(fHandler)
}

func MainPage(l logger.Logger, m json.Marshaler) http.Handler {
    fHandler :=  func(res http.ResponseWriter, req *http.Request) {

        var metricArr []metrics.Metric

        data, err := json.Marshal(m)

        if err != nil {
            l.Infoln(err)
            http.Error(res, "Error while loading main page", http.StatusNotFound)
        }

        err = json.Unmarshal(data, &metricArr)

        if err != nil {
            l.Infoln(err)
            http.Error(res, "Error while loading main page", http.StatusNotFound)
        }

        slices.SortStableFunc(metricArr, func(a,b metrics.Metric) int {
            return cmp.Compare(a.ID, b.ID)})

        var html bytes.Buffer

        html.WriteString(`<html>
                          <body>`)

        for _, v := range metricArr {
            if v.MType == "gauge" {
                html.WriteString(`<p>`)
                html.WriteString(v.ID)
                html.WriteString(" = ")
                html.WriteString(strconv.FormatFloat(*v.Value, 'f', -1, 64))
                html.WriteString(`</p>`)
            }
            if v.MType == "counter" {
                html.WriteString(`<p>`)
                html.WriteString(v.ID)
                html.WriteString(" = ")
                html.WriteString(strconv.FormatInt(*v.Delta, 10))
                html.WriteString(`</p>`)
            }
        }

        html.WriteString(`</html>
                          </body>`)


        res.Header().Add("Content-Type", "text/html")
        res.WriteHeader(http.StatusOK)
        res.Write(html.Bytes())
    }

    return http.HandlerFunc(fHandler)
}



func NotFound(res http.ResponseWriter, req *http.Request) {
    res.WriteHeader(http.StatusNotFound)
    res.Write([]byte("Not Found\n"))
}

func BadRequest(res http.ResponseWriter, req *http.Request) {
    res.WriteHeader(http.StatusBadRequest)
    res.Write([]byte("Bad Request\n"))
}
