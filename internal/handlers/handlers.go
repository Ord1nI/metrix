package handlers

import (
    "github.com/go-chi/chi/v5"

    "net/http"
    "strconv"
    "bytes"
    "io"
    "github.com/Ord1nI/metrix/internal/storage"
)


func UpdateGauge(s storage.Adder) func(res http.ResponseWriter, req *http.Request) {
    return func(res http.ResponseWriter, req *http.Request) {
        
        name := chi.URLParam(req, "name")
        v := chi.URLParam(req, "val")

        val,err := strconv.ParseFloat(v, 64)

        if err != nil {
            http.Error(res, "Incorect metric value", http.StatusBadRequest)
            return
        }

        s.AddGauge(name,val)
        res.WriteHeader(http.StatusOK)
    }
}

func UpdateCounter(s storage.Adder) func(res http.ResponseWriter, req *http.Request){
    return func(res http.ResponseWriter, req *http.Request) {

        name := chi.URLParam(req, "name")
        v := chi.URLParam(req, "val")
        
        val, err := strconv.ParseInt(v, 10, 64)

        if err != nil {
            http.Error(res, "Incorect metric value", http.StatusBadRequest)
            return
        }

        s.AddCounter(name, val)
        res.WriteHeader(http.StatusOK)
    }
}

func GetGauge(s storage.Getter) func(res http.ResponseWriter, req *http.Request) {
    return func(res http.ResponseWriter, req *http.Request) {
        name := chi.URLParam(req,"name")
        v, err := s.GetGauge(name)

        if err != nil {
            http.Error(res, "Unknown metric", http.StatusNotFound)
            return
        }

        res.WriteHeader(http.StatusOK)
        io.WriteString(res, strconv.FormatFloat(v, 'f', -1, 64))
        res.Write([]byte("\n"))
    }
}

func GetCounter(s storage.Getter) func(res http.ResponseWriter, req *http.Request){

    return func(res http.ResponseWriter, req *http.Request) {
        name := chi.URLParam(req,"name")
        v, err := s.GetCounter(name)

        if err != nil {
            http.Error(res, "Unknown metric", http.StatusNotFound)
            return
        }

        res.WriteHeader(http.StatusOK)
        io.WriteString(res, strconv.FormatInt(v, 10))
        res.Write([]byte("\n"))
    }

}
func GetAllMetrics(stor *storage.MemStorage) func(res http.ResponseWriter, req *http.Request) {
    return func(res http.ResponseWriter, req *http.Request) {
        var html bytes.Buffer
        html.WriteString(`<html>
                          <body>`)

        metrics := stor.GetAll()

        for _, i := range metrics {
            html.WriteString(`<p>`)
            html.WriteString(i.Name)
            html.WriteString(" = ")
            html.WriteString(strconv.FormatFloat(i.Val, 'f', -1, 64))
            html.WriteString(`</p>`)
        }
        html.WriteString(`</html>
                          </body>`)
        res.WriteHeader(http.StatusOK)
        res.Write(html.Bytes())
    }
}
func NotFound(res http.ResponseWriter, req *http.Request) {
    res.WriteHeader(http.StatusNotFound)
    res.Write([]byte("Not Found\n"))
}

func BadRequest(res http.ResponseWriter, req *http.Request) {
    res.WriteHeader(http.StatusBadRequest)
    res.Write([]byte("Bad Request\n"))
}
