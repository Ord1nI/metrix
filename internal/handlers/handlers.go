package handlers


import (
    "github.com/go-chi/chi/v5"

    "net/http"
    "strconv"
    "bytes"
    "sort"
    "io"
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

        s.AddGauge(name,storage.Gauge(val))
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

        s.AddCounter(name, storage.Counter(val))
        res.WriteHeader(http.StatusOK)
    }
    return http.HandlerFunc(fHandler)
}

func GetGauge(s storage.Getter) http.Handler {
    fHandler :=  func(res http.ResponseWriter, req *http.Request) {
        name := chi.URLParam(req,"name")
        v, err := s.GetGauge(name)

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
        v, err := s.GetCounter(name)

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

func GetAllMetrics(stor *storage.MemStorage) http.Handler {
    fHandler := func(res http.ResponseWriter, req *http.Request) {
        var html bytes.Buffer
        html.WriteString(`<html>
                          <body>`)

        GaugeNameArr := stor.GetGaugeNames()
        sort.Strings(GaugeNameArr)
        CounterNameArr := stor.GetCounterNames()
        sort.Strings(CounterNameArr)

        html.WriteString(`<b> GAUGE METRICS: </b>`)

        for _, i := range GaugeNameArr {
            html.WriteString(`<p>`)
            html.WriteString(i)
            html.WriteString(" = ")
            html.WriteString(strconv.FormatFloat(float64(stor.Gauge[i]), 'f', -1, 64))
            html.WriteString(`</p>`)
        }
        html.WriteString(`<b> COUNTER METRICS: </b>`)

        for _, i := range CounterNameArr {
            html.WriteString(`<p>`)
            html.WriteString(i)
            html.WriteString(" = ")
            html.WriteString(strconv.FormatInt(int64(stor.Counter[i]), 10))
            html.WriteString(`</p>`)
        }

        html.WriteString(`</html>
                          </body>`)
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
