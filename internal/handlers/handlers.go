package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgerrcode"

	"bytes"
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/Ord1nI/metrix/internal/logger"
	"github.com/Ord1nI/metrix/internal/repo"
	"github.com/Ord1nI/metrix/internal/repo/metrics"
)

var (
    ErrUpdate error = errors.New("error while updating")
    ErrGetting error = errors.New("error while getting")
    ErrNotJSON error = errors.New("not json request")
    ErrSQLuniqueViolation error = errors.New(pgerrcode.UniqueViolation)
    ErrSQLconnectionException error = errors.New(pgerrcode.ConnectionException)
)

type APIFunc func(http.ResponseWriter, *http.Request) error

func (a APIFunc) ServerHTTP(w http.ResponseWriter, r *http.Request) {
    err := a(w, r)
    if e, ok := err.(*HandlerError); ok {
        SendHandlerError(w, e)
        return 
    }
    SendInternalError(w)
}

type HandlerError struct {
	StatusCode int
	Err        error
}

func NewHandlerError(err error, status int) error {
	return &HandlerError{
		Err:        err,
		StatusCode: status,
	}
}

func (h *HandlerError) Error() string {
	return fmt.Sprintf("error: %s, with code %d", h.Err.Error(), h.StatusCode)
}

func SendInternalError(r http.ResponseWriter) {
    http.Error(r, "Internal server error", http.StatusInternalServerError)
}

func SendHandlerError(r http.ResponseWriter, err *HandlerError) {
    http.Error(r, err.Error(), err.StatusCode)
}


type APIHandler struct {
    l logger.Logger
    f func(http.ResponseWriter, *http.Request) error
    BackOffScedule []time.Duration
    BackOffErrors error
}

func NewAPIHandler(l logger.Logger, f APIFunc, backOffSchedule []time.Duration, errorList error) *APIHandler{
    return &APIHandler{
        l: l,
        f: f,
        BackOffScedule: backOffSchedule,
        BackOffErrors: errorList,
    }
}

func (a *APIHandler) ServeHTTP(res http.ResponseWriter, r *http.Request) {
    if err := a.f(res, r); err != nil {
        a.l.Infoln("Got Error: ", err)
        if apiErr, ok := err.(*HandlerError); ok {
            if a.BackOffScedule == nil {
                SendHandlerError(res, apiErr)
                return
            }

            if a.BackOffErrors == nil {
                SendHandlerError(res, apiErr)
                return
            }

            if !errors.Is(a.BackOffErrors, apiErr) {
                SendHandlerError(res, apiErr)
                return
            }

            if a.BackOff(res, r) {
                return
            }
            SendHandlerError(res, apiErr)
            return
        } else {
            SendInternalError(res)
            return
        }
    }
}

func (a *APIHandler) BackOff(res http.ResponseWriter, r *http.Request) bool{
    a.l.Infoln("Trying backoff handler")
	for _, backoff := range a.BackOffScedule {
        if err := a.f(res, r); err == nil {
			a.l.Infoln("Error successfuly recovered")
			return true
		}
		time.Sleep(backoff)
	}
    a.l.Infoln("Error wans't recover")
    return false
}

func UpdateGauge(s repo.Adder) APIFunc {
	fHandler := func(res http.ResponseWriter, req *http.Request) error {
		name := chi.URLParam(req, "name")
		v := chi.URLParam(req, "val")

		val, err := strconv.ParseFloat(v, 64)

		if err != nil {
			return NewHandlerError(errors.Join(err, ErrUpdate), http.StatusBadRequest)
		}

		s.Add(name, metrics.Gauge(val))
		res.WriteHeader(http.StatusOK)
		return nil
	}
	return APIFunc(fHandler)
}

func UpdateCounter(s repo.Adder) APIFunc {
	fHandler := func(res http.ResponseWriter, req *http.Request) error {

		name := chi.URLParam(req, "name")
		v := chi.URLParam(req, "val")

		val, err := strconv.ParseInt(v, 10, 64)

		if err != nil {
			return NewHandlerError(errors.Join(err, ErrUpdate), http.StatusBadRequest)
		}

		s.Add(name, metrics.Counter(val))
		res.WriteHeader(http.StatusOK)
		return nil
	}
	return APIFunc(fHandler)
}

func GetGauge(s repo.Getter) APIFunc {
	fHandler := func(res http.ResponseWriter, req *http.Request) error {
		name := chi.URLParam(req, "name")
		var v metrics.Gauge
		err := s.Get(name, &v)

		if err != nil {
			return NewHandlerError(errors.Join(err, ErrGetting), http.StatusNotFound)
		}

		res.WriteHeader(http.StatusOK)
		io.WriteString(res, strconv.FormatFloat(float64(v), 'f', -1, 64))
		res.Write([]byte("\n"))
		return nil
	}
	return APIFunc(fHandler)
}

func GetCounter(s repo.Getter) APIFunc {

	fHandler := func(res http.ResponseWriter, req *http.Request) error {
		name := chi.URLParam(req, "name")
		var v metrics.Counter
		err := s.Get(name, &v)

		if err != nil {
			return NewHandlerError(errors.Join(err, ErrGetting), http.StatusNotFound)
		}

		res.WriteHeader(http.StatusOK)
		io.WriteString(res, strconv.FormatInt(int64(v), 10))
		res.Write([]byte("\n"))
		return nil
	}
	return APIFunc(fHandler)
}

func MainPage(m json.Marshaler) APIFunc {
	fHandler := func(res http.ResponseWriter, req *http.Request) error {

		var metricArr []metrics.Metric

		data, err := json.Marshal(m)

		if err != nil {
			return NewHandlerError(errors.Join(err, errors.New("Error while loading main page")), http.StatusBadRequest)
		}

		err = json.Unmarshal(data, &metricArr)

		if err != nil {
			return NewHandlerError(errors.Join(err, errors.New("Error while loading main page")), http.StatusBadRequest)
		}

		slices.SortStableFunc(metricArr, func(a, b metrics.Metric) int {
			return cmp.Compare(a.ID, b.ID)
		})

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
