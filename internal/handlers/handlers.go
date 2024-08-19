package handlers

import (
	"github.com/go-chi/chi/v5"
    "github.com/jackc/pgerrcode"

	"bytes"
	"cmp"
	"encoding/json"
    "fmt"
	"io"
    "time"
	"net/http"
	"slices"
    "errors"
	"strconv"

	"github.com/Ord1nI/metrix/internal/repo/metrics"
	"github.com/Ord1nI/metrix/internal/repo"
	"github.com/Ord1nI/metrix/internal/logger"
)
var errList = errors.Join(errors.New(pgerrcode.UniqueViolation),errors.New(pgerrcode.ConnectionException))

type HandlerError struct {
    StatusCode int
    Err error
    Retry bool
}

func (h HandlerError) Error()string {
    return fmt.Sprintf("error: %s, with code %d",h.Err.Error(), h.StatusCode)
}

func NewHandlerError(err error, status int) error{
    return HandlerError {
        Err:err,
        StatusCode: status,
        Retry:errors.Is(err, errList),
    }
}

type APIFunc func(http.ResponseWriter, *http.Request) error

func Make(l logger.Logger, f APIFunc, BackoffSchedule []time.Duration) http.Handler{
    fun := func(w http.ResponseWriter, r *http.Request) {
        if err := f(w,r); err != nil {
            l.Infoln(err)
            if apiErr, ok := err.(HandlerError); ok {
                if apiErr.Retry {
                    for _, backoff := range BackoffSchedule {
                        if err = f(w,r); err == nil {
                            l.Infoln("Error successfuly recovered")
                            return
                        }
                        time.Sleep(backoff)
                    }
                }
                http.Error(w, apiErr.Error(),apiErr.StatusCode)
                return
            } else {
                http.Error(w, "internal server error", http.StatusInternalServerError)
            }
        }
    }
    return http.HandlerFunc(fun)
}

func UpdateGauge(s repo.Adder) APIFunc {
    fHandler := func(res http.ResponseWriter, req *http.Request) error {
        name := chi.URLParam(req, "name")
        v := chi.URLParam(req, "val")

        val,err := strconv.ParseFloat(v, 64)

       if err != nil {
            return NewHandlerError(err,http.StatusBadRequest)
        }

        s.Add(name, metrics.Gauge(val))
        res.WriteHeader(http.StatusOK)
        return nil
    }
    return APIFunc(fHandler)
}

func UpdateCounter(s repo.Adder) APIFunc{
    fHandler := func(res http.ResponseWriter, req *http.Request) error{

        name := chi.URLParam(req, "name")
        v := chi.URLParam(req, "val")
        
        val, err := strconv.ParseInt(v, 10, 64)

        if err != nil {
            return NewHandlerError(err,http.StatusBadRequest)
        }

        s.Add(name, metrics.Counter(val))
        res.WriteHeader(http.StatusOK)
        return nil
    }
    return APIFunc(fHandler)
}

func GetGauge(s repo.Getter) APIFunc {
    fHandler :=  func(res http.ResponseWriter, req *http.Request) error {
        name := chi.URLParam(req,"name")
        var v metrics.Gauge
        err := s.Get(name, &v)

        if err != nil {
            return NewHandlerError(err,http.StatusNotFound)
        }

        res.WriteHeader(http.StatusOK)
        io.WriteString(res, strconv.FormatFloat(float64(v), 'f', -1, 64))
        res.Write([]byte("\n"))
        return nil
    }
    return APIFunc(fHandler)
}

func GetCounter(s repo.Getter) APIFunc {

    fHandler :=  func(res http.ResponseWriter, req *http.Request) error {
        name := chi.URLParam(req,"name")
        var v metrics.Counter
        err := s.Get(name, &v)

        if err != nil {
            return NewHandlerError(err,http.StatusNotFound)
        }

        res.WriteHeader(http.StatusOK)
        io.WriteString(res, strconv.FormatInt(int64(v), 10))
        res.Write([]byte("\n"))
        return nil
    }
    return APIFunc(fHandler)
}

func MainPage(m json.Marshaler) APIFunc {
    fHandler :=  func(res http.ResponseWriter, req *http.Request)error {

        var metricArr []metrics.Metric

        data, err := json.Marshal(m)

        if err != nil {
            return NewHandlerError(err,http.StatusBadRequest)
        }

        err = json.Unmarshal(data, &metricArr)

        if err != nil {
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
        return nil
    }

    return APIFunc(fHandler)
}



func NotFound(res http.ResponseWriter, req *http.Request) {
    res.WriteHeader(http.StatusNotFound)
    res.Write([]byte("Not Found\n"))
}

func BadRequest(res http.ResponseWriter, req *http.Request) {
    res.WriteHeader(http.StatusBadRequest)
    res.Write([]byte("Bad Request\n"))
}
