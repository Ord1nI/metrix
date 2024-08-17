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

	"github.com/Ord1nI/metrix/internal/storage"
)


func UpdateGauge(s storage.Adder) http.Handler {
    fHandler := func(res http.ResponseWriter, req *http.Request) {
        name := chi.URLParam(req, "name")
        v := chi.URLParam(req, "val")

        val,err := strconv.ParseFloat(v, 64)

        if err != nil {
            http.Error(res, "Incorect metric value", http.StatusBadRequest)
            return
        }

        s.Add(name, storage.Gauge(val))
        res.WriteHeader(http.StatusOK)
    }
    return http.HandlerFunc(fHandler)
}

func UpdateCounter(s storage.Adder) http.Handler{
    fHandler := func(res http.ResponseWriter, req *http.Request) {

        name := chi.URLParam(req, "name")
        v := chi.URLParam(req, "val")
        
        val, err := strconv.ParseInt(v, 10, 64)

        if err != nil {
            http.Error(res, "Incorect metric value", http.StatusBadRequest)
            return
        }

        s.Add(name, storage.Counter(val))
        res.WriteHeader(http.StatusOK)
    }
    return http.HandlerFunc(fHandler)
}

func GetGauge(s storage.Getter) http.Handler {
    fHandler :=  func(res http.ResponseWriter, req *http.Request) {
        name := chi.URLParam(req,"name")
        var v storage.Gauge
        err := s.Get(name, &v)

        if err != nil {
            http.Error(res, "Unknown metric", http.StatusNotFound)
            return
        }

        res.WriteHeader(http.StatusOK)
        io.WriteString(res, strconv.FormatFloat(float64(v), 'f', -1, 64))
        res.Write([]byte("\n"))
    }
    return http.HandlerFunc(fHandler)
}

func GetCounter(s storage.Getter) http.Handler {

    fHandler :=  func(res http.ResponseWriter, req *http.Request) {
        name := chi.URLParam(req,"name")
        var v storage.Counter
        err := s.Get(name, &v)

        if err != nil {
            http.Error(res, "Unknown metric", http.StatusNotFound)
            return
        }

        res.WriteHeader(http.StatusOK)
        io.WriteString(res, strconv.FormatInt(int64(v), 10))
        res.Write([]byte("\n"))
    }
    return http.HandlerFunc(fHandler)
}
func MainPage(m json.Marshaler) http.Handler {
    fHandler :=  func(res http.ResponseWriter, req *http.Request) {

        var metricArr []storage.Metric

        data, err := json.Marshal(m)

        if err != nil {
            http.Error(res, "error", http.StatusNotFound)
        }

        err = json.Unmarshal(data, &metricArr)

        if err != nil {
            http.Error(res, "couldn't unmarshal", http.StatusNotFound)
        }

        slices.SortStableFunc(metricArr, func(a,b storage.Metric) int {
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
